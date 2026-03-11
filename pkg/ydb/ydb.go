package ydb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/grafana/grafana/pkg/services/sqlstore/migrator"
	"github.com/grafana/grafana/pkg/services/sqlstore/session"
	"github.com/grafana/grafana/pkg/storage/unified/sql/sqltemplate"
	"github.com/grafana/grafana/pkg/ydb/filters"
	"github.com/grafana/grafana/pkg/ydb/snapshot"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/pkg/xerrors"
	"github.com/ydb-platform/ydb-go-sdk/v3/pkg/xslices"
	"github.com/ydb-platform/ydb-go-sdk/v3/retry"
	"github.com/ydb-platform/ydb-go-sdk/v3/scheme"
	"github.com/ydb-platform/ydb-go-sdk/v3/sugar"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/options"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
	yc "github.com/ydb-platform/ydb-go-yc-metadata"
	yql "github.com/ydb-platform/yql-parsers/go"
	"golang.org/x/exp/slices"

	"github.com/grafana/grafana/pkg/util/xorm"
	"github.com/grafana/grafana/pkg/util/xorm/core"
)

func init() {
	ydb.RegisterDsnParser(func(dsn string) (opts []ydb.Option, _ error) {
		uri, err := url.Parse(dsn)
		if err != nil {
			return opts, nil
		}

		if ycAuth := uri.Query().Get("yc_auth"); ycAuth == "metadata_credentials" {
			opts = append(opts,
				ydb.WithCredentials(yc.NewInstanceServiceAccount()),
				yc.WithInternalCA(),
			)
		}

		return opts, nil
	})
}

// from https://github.com/ydb-platform/ydb/blob/main/ydb/library/yql/sql/v1/SQLv1.g.in#L1117
var (
	ydbReservedWords = xslices.Map(yql.NewYQLLexer(nil).GetSymbolicNames(), func(reservedWord string) string {
		return reservedWord
	})
)

func toYQLDataType(t string, isAutoIncrement bool) string {
	switch t {
	case core.Bool, core.Boolean:
		return types.TypeBool.Yql()
	case core.TinyInt:
		return types.TypeInt64.Yql()
	case core.Int, core.Integer:
		if isAutoIncrement {
			return "Serial"
		}
		return types.TypeInt64.Yql()
	case core.SmallInt, core.MediumInt, core.BigInt:
		if isAutoIncrement {
			return "Serial"
		}
		return types.TypeInt64.Yql()
	case core.Float:
		return types.TypeDouble.Yql()
	case core.Double:
		return types.TypeDouble.Yql()
	case core.Blob, core.LongBlob, core.MediumBlob, core.TinyBlob, core.VarBinary, core.Binary:
		return types.TypeBytes.Yql()
	case core.Json:
		return types.TypeJSONDocument.Yql()
	case core.Varchar, core.NVarchar, core.Char, core.NChar,
		core.MediumText, core.LongText, core.Text, core.NText, core.TinyText:
		return types.TypeText.Yql()
	case core.TimeStamp, core.Time, core.Date, core.DateTime:
		return types.TypeText.Yql()
	case core.Serial, core.BigSerial:
		return "Serial"
	default:
		return t
	}
}

func yqlToSQLType(yqlType string) core.SQLType {
	switch yqlType {
	case types.TypeBool.Yql():
		return core.SQLType{Name: core.Bool, DefaultLength: 0, DefaultLength2: 0}
	case types.TypeInt8.Yql(), types.TypeInt16.Yql(), types.TypeInt32.Yql(), types.TypeInt64.Yql(),
		types.TypeUint8.Yql(), types.TypeUint16.Yql(), types.TypeUint32.Yql(), types.TypeUint64.Yql():
		return core.SQLType{Name: core.BigInt, DefaultLength: 0, DefaultLength2: 0}
	case types.TypeDouble.Yql(), types.TypeFloat.Yql():
		return core.SQLType{Name: core.Double, DefaultLength: 0, DefaultLength2: 0}
	case types.TypeBytes.Yql(), types.TypeString.Yql():
		return core.SQLType{Name: core.Blob, DefaultLength: 0, DefaultLength2: 0}
	case types.TypeJSON.Yql(), types.TypeJSONDocument.Yql():
		return core.SQLType{Name: core.Json, DefaultLength: 0, DefaultLength2: 0}
	case types.TypeText.Yql(), types.TypeUTF8.Yql():
		return core.SQLType{Name: core.Varchar, DefaultLength: 0, DefaultLength2: 0}
	case types.TypeTimestamp.Yql(), types.TypeTimestamp64.Yql():
		return core.SQLType{Name: core.TimeStamp, DefaultLength: 0, DefaultLength2: 0}
	case types.TypeDatetime.Yql(), types.TypeDatetime64.Yql():
		return core.SQLType{Name: core.DateTime, DefaultLength: 0, DefaultLength2: 0}
	default:
		return core.SQLType{Name: yqlType}
	}
}

var (
	_ core.Driver                 = (*YDB)(nil)
	_ core.Dialect                = (*YDB)(nil)
	_ core.DialectWithReturningID = (*YDB)(nil)
	_ migrator.Dialect            = (*YDB)(nil)
	_ sqltemplate.Dialect         = (*YDB)(nil)
)

const ydbDriverName = "ydb"

func NewYDB() *YDB {
	return &YDB{
		log: &Logger{},
	}
}

type YDB struct {
	core.Base

	db  *ydb.Driver
	log *Logger

	tableParams map[string]string // TODO: maybe remove
}

func (d *YDB) DialectName() string {
	return ydbDriverName
}

func (d *YDB) Ident(s string) (string, error) {
	return d.Quote(s), nil
}

func (d *YDB) ArgPlaceholder(argNum int) string {
	return "$p" + strconv.Itoa(argNum+1)
}

func (d *YDB) SelectFor(s ...string) (string, error) {
	//TODO implement me
	panic("implement me")
}

func (d *YDB) CurrentEpoch() string {
	//TODO implement me
	panic("implement me")
}

// DSN format: https://github.com/ydb-platform/ydb-go-sdk/blob/a804c31be0d3c44dfd7b21ed49d863619217b11d/connection.go#L339
func (d *YDB) Parse(driverName, dataSourceName string) (*core.Uri, error) {
	if driverName != ydbDriverName {
		return nil, xerrors.WithStackTrace(fmt.Errorf("wrong driver name: %q", driverName))
	}

	info := &core.Uri{DbType: ydbDriverName}

	uri, err := url.Parse(dataSourceName)
	if err != nil {
		return nil, xerrors.WithStackTrace(fmt.Errorf("failed on parse data source %v", dataSourceName))
	}

	const (
		secure   = "grpcs"
		insecure = "grpc"
	)

	if uri.Scheme != secure && uri.Scheme != insecure {
		return nil, xerrors.WithStackTrace(fmt.Errorf("unsupported scheme %v", uri.Scheme))
	}

	info.Host = uri.Host
	if spl := strings.Split(uri.Host, ":"); len(spl) > 1 {
		info.Host = spl[0]
		info.Port = spl[1]
	}

	info.DbName = uri.Path
	if info.DbName == "" {
		return nil, xerrors.WithStackTrace(errors.New("database path can not be empty"))
	}

	if uri.User != nil {
		info.Passwd, _ = uri.User.Password()
		info.User = uri.User.Username()
	}

	return info, nil
}

func (d *YDB) Quote(s string) string {
	return "`" + s + "`"
}

// DriverName returns the dialect driver name so that migrator gets "ydb" even when Init was not called.
func (d *YDB) DriverName() string {
	return ydbDriverName
}

// RecursiveQueriesAreSupported implements migrator.DialectRecursiveCTE (optional).
func (d *YDB) RecursiveQueriesAreSupported() (ok bool, err error) {
	end := d.log.Start("RecursiveQueriesAreSupported")
	defer func() { end.End(ok, err) }()
	return false, nil
}

func (d *YDB) SQLType(col *migrator.Column) (res string) {
	end := d.log.Start("SQLType", col)
	defer func() { end.End(res) }()
	return toYQLDataType(col.Type, col.IsAutoIncrement)
}

func (d *YDB) LikeOperator(column string, wildcardBefore bool, pattern string, wildcardAfter bool) (sql string, param string) {
	end := d.log.Start("LikeOperator", column, wildcardBefore, pattern, wildcardAfter)
	defer func() { end.End(sql, param) }()
	param = pattern
	if wildcardBefore {
		param = "%" + param
	}
	if wildcardAfter {
		param = param + "%"
	}
	return column + " ILIKE ?", param
}

func (d *YDB) Default(col *migrator.Column) (res string) {
	end := d.log.Start("Default", col)
	defer func() { end.End(res) }()
	if col.Type == migrator.DB_Bool {
		bl, err := strconv.ParseBool(col.Default)
		if err != nil {
			panic(fmt.Errorf("failed to create default value for column '%s': invalid boolean default value '%s'", col.Name, col.Default))
		}
		return d.BooleanStr(bl)
	}
	if col.Type == migrator.DB_NVarchar {
		return `"` + col.Default + `"`
	}
	return col.Default
}

func (d *YDB) BooleanValue(b bool) (res any) {
	end := d.log.Start("BooleanValue", b)
	defer func() { end.End(res) }()
	return b
}

func (d *YDB) BooleanStr(b bool) (res string) {
	end := d.log.Start("BooleanStr", b)
	defer func() { end.End(res) }()
	return strconv.FormatBool(b)
}

func (d *YDB) DateTimeFunc(s string) (res string) {
	end := d.log.Start("DateTimeFunc", s)
	defer func() { end.End(res) }()
	return s
}

func (d *YDB) BatchSize() (res int) {
	end := d.log.Start("BatchSize")
	defer func() { end.End(res) }()
	return 1000
}

func (d *YDB) UnionDistinct() (res string) {
	end := d.log.Start("UnionDistinct")
	defer func() { end.End(res) }()
	return "UNION"
}

func (d *YDB) UnionAll() (res string) {
	end := d.log.Start("UnionAll")
	defer func() { end.End(res) }()
	return "UNION ALL"
}

func (d *YDB) OrderBy(order string) (res string) {
	end := d.log.Start("OrderBy", order)
	defer func() { end.End(res) }()
	return strings.ReplaceAll(order, "dashboard.title", "title")
}

func (d *YDB) CreateIndexSQL(tableName string, index *migrator.Index) (res string) {
	end := d.log.Start("CreateIndexSQL", tableName, index)
	defer func() { end.End(res) }()
	indexName := d.Quote(index.XName(tableName))
	tableName = d.Quote(tableName)
	colsIndex := make([]string, len(index.Cols))
	for i := range index.Cols {
		colsIndex[i] = d.Quote(index.Cols[i])
	}
	indexOn := strings.Join(colsIndex, ",")
	var buf strings.Builder
	buf.WriteString("ALTER TABLE ")
	buf.WriteString(tableName)
	buf.WriteString(" ADD INDEX ")
	buf.WriteString(indexName)
	buf.WriteString(" GLOBAL ON ( ")
	buf.WriteString(indexOn)
	buf.WriteString(" );")
	return buf.String()
}

func (d *YDB) CreateTableSQL(table *migrator.Table) (res string) {
	end := d.log.Start("CreateTableSQL", table)
	defer func() { end.End(res) }()
	var sql strings.Builder
	sql.WriteString("CREATE TABLE IF NOT EXISTS ")
	sql.WriteString(d.Quote(table.Name))
	sql.WriteString(" (\n")

	if len(table.PrimaryKeys) == 0 {
		table.Columns = append([]*migrator.Column{
			{Name: "id", Type: migrator.DB_Serial, IsPrimaryKey: true, IsAutoIncrement: true},
		}, table.Columns...)
		table.PrimaryKeys = []string{"id"}
	}

	for i, col := range table.Columns {
		if i > 0 {
			sql.WriteString(",\n")
		}
		sql.WriteString("\t")
		sql.WriteString(col.StringNoPk(d))
	}

	sql.WriteString(",\n\tPRIMARY KEY (")
	for i, col := range table.PrimaryKeys {
		if i > 0 {
			sql.WriteString(",")
		}
		sql.WriteString(d.Quote(col))
	}
	sql.WriteString(")")

	if len(table.Indices) > 0 {
		for _, idx := range table.Indices {
			if idx.Type == migrator.UniqueIndex {
				colNames := make([]string, len(table.Columns))
				for i, c := range table.Columns {
					colNames[i] = c.Name
				}
				if len(idx.Cols) == 1 && slices.Contains(colNames, idx.Name) && len(table.PrimaryKeys) == 1 && slices.Contains(table.PrimaryKeys, idx.Name) {
					continue
				}
			}
			if idx.Name == "" {
				idx.Name = table.Name + "_" + strings.Join(idx.Cols, "_")
				if idx.Type == migrator.UniqueIndex {
					idx.Name = "UQE_" + idx.Name
				}
			}
			sql.WriteString(",\n\tINDEX ")
			sql.WriteString(d.Quote(idx.Name))
			sql.WriteString(" GLOBAL")
			if idx.Type == migrator.UniqueIndex {
				sql.WriteString(" UNIQUE")
			}
			sql.WriteString(" SYNC ON (")
			for i, col := range idx.Cols {
				if i > 0 {
					sql.WriteString(",")
				}
				sql.WriteString(d.Quote(col))
			}
			sql.WriteString(")")
		}
	}

	sql.WriteString("\n);")
	return sql.String()
}

func (d *YDB) AddColumnSQL(tableName string, col *migrator.Column) (res string) {
	end := d.log.Start("AddColumnSQL", tableName, col)
	defer func() { end.End(res) }()
	col.Nullable = true
	col.Default = ""
	var buf strings.Builder
	buf.WriteString("ALTER TABLE ")
	buf.WriteString(d.Quote(tableName))
	buf.WriteString(" ADD COLUMN ")
	buf.WriteString(d.ColStringNoPk(col))
	return buf.String()
}

func (d *YDB) CopyTableData(sourceTable string, targetTable string, sourceCols []string, targetCols []string) (res string) {
	end := d.log.Start("CopyTableData", sourceTable, targetTable, sourceCols, targetCols)
	defer func() { end.End(res) }()
	sourceColsSQL := ""
	for i, col := range sourceCols {
		if i > 0 {
			sourceColsSQL += ", "
		}
		sourceColsSQL += d.Quote(col)
	}
	targetColsSQL := ""
	for i, col := range targetCols {
		if i > 0 {
			targetColsSQL += ", "
		}
		targetColsSQL += d.Quote(col)
	}
	var buf strings.Builder
	buf.WriteString("INSERT INTO ")
	buf.WriteString(d.Quote(targetTable))
	buf.WriteString(" (")
	buf.WriteString(targetColsSQL)
	buf.WriteString(") SELECT ")
	buf.WriteString(sourceColsSQL)
	buf.WriteString(" FROM ")
	buf.WriteString(d.Quote(sourceTable))
	return buf.String()
}

func (d *YDB) DropTable(tableName string) (res string) {
	end := d.log.Start("DropTable", tableName)
	defer func() { end.End(res) }()
	return "DROP TABLE IF EXISTS " + d.Quote(tableName)
}

func (d *YDB) DropIndexSQL(tableName string, index *migrator.Index) (res string) {
	end := d.log.Start("DropIndexSQL", tableName, index)
	defer func() { end.End(res) }()
	var buf strings.Builder
	buf.WriteString("ALTER TABLE ")
	buf.WriteString(d.Quote(tableName))
	buf.WriteString(" DROP INDEX ")
	buf.WriteString(d.Quote(index.XName(tableName)))
	return buf.String()
}

func (d *YDB) RenameTable(oldName string, newName string) (res string) {
	end := d.log.Start("RenameTable", oldName, newName)
	defer func() { end.End(res) }()
	var buf strings.Builder
	buf.WriteString("ALTER TABLE ")
	buf.WriteString(d.Quote(oldName))
	buf.WriteString(" RENAME TO ")
	buf.WriteString(d.Quote(newName))
	return buf.String()
}

func (d *YDB) RenameColumn(table migrator.Table, column *migrator.Column, newName string) (res string) {
	end := d.log.Start("RenameColumn", table, column, newName)
	defer func() { end.End(res) }()
	oldName := column.Name
	column.Name = newName
	sql := d.AddColumnSQL(table.Name, column) + ";"
	if !column.IsPrimaryKey {
		column.Name = oldName
		sql += "ALTER TABLE " + d.Quote(table.Name) + " DROP COLUMN " + d.Quote(column.Name)
	}
	return sql
}

func (d *YDB) UpdateTableSQL(tableName string, columns []*migrator.Column) (res string) {
	end := d.log.Start("UpdateTableSQL", tableName, columns)
	defer func() { end.End(res) }()
	return "-- NOT REQUIRED"
}

func (d *YDB) IndexCheckSQL(tableName, indexName string) (sql string, args []any) {
	end := d.log.Start("IndexCheckSQL", tableName, indexName)
	defer func() { end.End(sql, args) }()

	return "SELECT Path FROM `.sys/partition_stats` where Path LIKE '%/'" +
		" || ? || '/' || ? || '/indexImplTable'", []any{tableName, indexName}
}

func (d *YDB) ColumnCheckSQL(tableName, columnName string) (sql string, args []any) {
	end := d.log.Start("ColumnCheckSQL", tableName, columnName)
	defer func() { end.End(sql, args) }()
	return "", nil
}

func (d *YDB) UpsertSQL(tableName string, keyCols, updateCols []string) (res string) {
	end := d.log.Start("UpsertSQL", tableName, keyCols, updateCols)
	defer func() { end.End(res) }()
	str, _ := d.UpsertMultipleSQL(tableName, keyCols, updateCols, 1)
	return str
}

func (d *YDB) UpsertMultipleSQL(tableName string, keyCols, updateCols []string, count int) (res string, err error) {
	end := d.log.Start("UpsertMultipleSQL", tableName, keyCols, updateCols, count)
	defer func() { end.End(res, err) }()
	if count < 1 {
		return "", fmt.Errorf("upsert statement must have count >= 1. Got %v", count)
	}
	columnsStr := strings.Builder{}
	valuesStr := strings.Builder{}
	for i, c := range updateCols {
		if i > 0 {
			columnsStr.WriteString(", ")
			valuesStr.WriteString(", ")
		}
		columnsStr.WriteString(d.Quote(c))
		valuesStr.WriteString("?")
	}
	var result strings.Builder
	result.WriteString("UPSERT INTO ")
	result.WriteString(d.Quote(tableName))
	result.WriteString(" (")
	result.WriteString(columnsStr.String())
	result.WriteString(") VALUES (")
	result.WriteString(valuesStr.String())
	result.WriteString(");")
	return result.String(), nil
}

func (d *YDB) ColString(column *migrator.Column) (res string) {
	end := d.log.Start("ColString", column)
	defer func() { end.End(res) }()
	sql := d.Quote(column.Name) + " "
	sql += d.SQLType(column) + " "
	if column.IsPrimaryKey {
		sql += "PRIMARY KEY " + d.AutoIncrStr() + " "
	}
	if column.Nullable {
		sql += "NULL "
	} else {
		sql += "NOT NULL "
	}
	if column.Default != "" {
		sql += "DEFAULT " + d.Default(column) + " "
	}
	return sql
}

func (d *YDB) ColStringNoPk(column *migrator.Column) (res string) {
	end := d.log.Start("ColStringNoPk", column)
	defer func() { end.End(res) }()
	sql := d.Quote(column.Name) + " "
	sql += d.SQLType(column)
	if column.Default != "" {
		sql += " DEFAULT " + d.Default(column) + " "
	}
	return sql
}

func (d *YDB) Limit(limit int64) (res string) {
	end := d.log.Start("Limit", limit)
	defer func() { end.End(res) }()
	return " LIMIT " + strconv.FormatInt(limit, 10)
}

func (d *YDB) LimitOffset(limit int64, offset int64) (res string) {
	end := d.log.Start("LimitOffset", limit, offset)
	defer func() { end.End(res) }()
	return " LIMIT " + strconv.FormatInt(limit, 10) + " OFFSET " + strconv.FormatInt(offset, 10)
}

func (d *YDB) PreInsertId(table string, sess *xorm.Session) (err error) {
	end := d.log.Start("PreInsertId", table, sess)
	defer func() { end.End(err) }()
	return nil
}

func (d *YDB) PostInsertId(table string, sess *xorm.Session) (err error) {
	end := d.log.Start("PostInsertId", table, sess)
	defer func() { end.End(err) }()
	return nil
}

func (d *YDB) CleanDB(engine *xorm.Engine) (err error) {
	end := d.log.Start("CleanDB", engine)
	defer func() { end.End(err) }()
	tables, err := engine.Dialect().GetTables()
	if err != nil {
		return err
	}
	sess := engine.NewSession()
	defer sess.Close()
	for _, table := range tables {
		if table.Name == "" {
			continue
		}
		dropSQL := "DROP TABLE IF EXISTS " + d.Quote(table.Name) + ";"
		if _, err := sess.Exec(dropSQL); err != nil {
			return fmt.Errorf("failed to drop table %q: %w", table.Name, err)
		}
	}
	return nil
}

func (d *YDB) TruncateDBTables(engine *xorm.Engine) (err error) {
	end := d.log.Start("TruncateDBTables", engine)
	defer func() { end.End(err) }()
	tables, err := engine.Dialect().GetTables()
	if err != nil {
		return err
	}
	sess := engine.NewSession()
	defer sess.Close()

	for _, table := range tables {
		switch table.Name {
		case "":
			continue
		case "migration_log":
			continue
		case "dashboard_acl":
			deleteSQL := "DELETE FROM " + d.Quote(table.Name) + " WHERE dashboard_id != -1 AND org_id != -1;"
			if _, err := sess.Exec(deleteSQL); err != nil {
				return fmt.Errorf("failed to truncate table %q: %w", table.Name, err)
			}
			continue
		default:
			deleteSQL := "DELETE FROM " + d.Quote(table.Name) + ";"
			err := retry.Do(context.Background(), engine.DB().DB, func(ctx context.Context, cc *sql.Conn) error {
				_, err := sess.Exec(deleteSQL)
				return err
			})
			if err != nil {
				if ydb.IsOperationErrorSchemeError(err) {
					continue
				}
				return fmt.Errorf("failed to truncate table %q: %w", table.Name, err)
			}
		}
	}
	return nil
}

func (d *YDB) CreateDatabaseFromSnapshot(ctx context.Context, engine *xorm.Engine, migrationLogTableName string) (err error) {
	end := d.log.Start("CreateDatabaseFromSnapshot", ctx, engine, migrationLogTableName)
	defer func() { end.End(err) }()

	if err := snapshot.Apply(ctx, d.db, migrationLogTableName); err != nil {
		return xerrors.WithStackTrace(err)
	}

	return nil
}

func (d *YDB) IsUniqueConstraintViolation(err error) (ok bool) {
	end := d.log.Start("IsUniqueConstraintViolation", err)
	defer func() { end.End(ok) }()
	return ydb.IsOperationErrorSchemeError(err)
}

func (d *YDB) ErrorMessage(err error) (res string) {
	end := d.log.Start("ErrorMessage", err)
	defer func() { end.End(res) }()
	if err != nil {
		return err.Error()
	}
	return ""
}

func (d *YDB) IsDeadlock(err error) (ok bool) {
	end := d.log.Start("IsDeadlock", err)
	defer func() { end.End(ok) }()
	return false
}

func (d *YDB) Lock(cfg migrator.LockCfg) (err error) {
	end := d.log.Start("Lock", cfg)
	defer func() { end.End(err) }()
	return nil
}

func (d *YDB) Unlock(cfg migrator.LockCfg) (err error) {
	end := d.log.Start("Unlock", cfg)
	defer func() { end.End(err) }()
	return nil
}

func (d *YDB) GetDBName(s string) (name string, err error) {
	end := d.log.Start("GetDBName", s)
	defer func() { end.End(name, err) }()
	uri, err := url.Parse(s)
	if err != nil {
		return "", fmt.Errorf("failed on parse data source %v", s)
	}
	return uri.Path, nil
}

func (d *YDB) InsertQuery(tableName string, row map[string]any) (query string, args []any, err error) {
	end := d.log.Start("InsertQuery", tableName, row)
	defer func() { end.End(query, args, err) }()
	if len(row) < 1 {
		return "", nil, fmt.Errorf("no columns provided")
	}
	keys := make([]string, 0, len(row))
	for col := range row {
		keys = append(keys, col)
	}
	slices.Sort(keys)
	cols := make([]string, 0, len(row))
	vals := make([]any, 0, len(row))
	for _, col := range keys {
		cols = append(cols, d.Quote(col))
		vals = append(vals, row[col])
	}
	placeholders := strings.Repeat("?, ", len(row)-1) + "?"
	var buf strings.Builder
	buf.WriteString("INSERT INTO ")
	buf.WriteString(d.Quote(tableName))
	buf.WriteString(" (")
	buf.WriteString(strings.Join(cols, ", "))
	buf.WriteString(") VALUES (")
	buf.WriteString(placeholders)
	buf.WriteString(")")
	return buf.String(), vals, nil
}

func (d *YDB) UpdateQuery(tableName string, row map[string]any, where map[string]any) (query string, args []any, err error) {
	end := d.log.Start("UpdateQuery", tableName, row, where)
	defer func() { end.End(query, args, err) }()
	if len(row) < 1 {
		return "", nil, fmt.Errorf("no columns provided")
	}
	if len(where) < 1 {
		return "", nil, fmt.Errorf("no where clause provided")
	}
	keys := make([]string, 0, len(row))
	for col := range row {
		keys = append(keys, col)
	}
	slices.Sort(keys)
	cols := make([]string, 0, len(row))
	vals := make([]any, 0, len(row)+len(where))
	for _, col := range keys {
		cols = append(cols, d.Quote(col)+"=?")
		vals = append(vals, row[col])
	}
	whereKeys := make([]string, 0, len(where))
	for col := range where {
		whereKeys = append(whereKeys, col)
	}
	slices.Sort(whereKeys)
	whereCols := make([]string, 0, len(where))
	for _, col := range whereKeys {
		whereCols = append(whereCols, d.Quote(col)+"=?")
		vals = append(vals, where[col])
	}
	var buf strings.Builder
	buf.WriteString("UPDATE ")
	buf.WriteString(d.Quote(tableName))
	buf.WriteString(" SET ")
	buf.WriteString(strings.Join(cols, ", "))
	buf.WriteString(" WHERE ")
	buf.WriteString(strings.Join(whereCols, " AND "))
	return buf.String(), vals, nil
}

func (d *YDB) Insert(ctx context.Context, tx *session.SessionTx, tableName string, row map[string]any) (err error) {
	end := d.log.Start("Insert", ctx, tx, tableName, row)
	defer func() { end.End(err) }()
	query, args, err := d.InsertQuery(tableName, row)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, query, args...)
	return err
}

func (d *YDB) Update(ctx context.Context, tx *session.SessionTx, tableName string, row map[string]any, where map[string]any) (err error) {
	end := d.log.Start("Update", ctx, tx, tableName, row, where)
	defer func() { end.End(err) }()
	query, args, err := d.UpdateQuery(tableName, row, where)
	if err != nil {
		return err
	}
	_, err = tx.Exec(ctx, query, args...)
	return err
}

func (d *YDB) Concat(s ...string) (res string) {
	end := d.log.Start("Concat", s)
	defer func() { end.End(res) }()
	return strings.Join(s, " || ")
}

func (d *YDB) WithReturningID() {
	end := d.log.Start("WithReturningID")
	defer func() { end.End() }()
}

func (d *YDB) Init(db *core.DB, uri *core.Uri, driverName, dataSource string) error {
	if driverName != ydbDriverName {
		return xerrors.WithStackTrace(fmt.Errorf("wrong driver name: %q", driverName))
	}

	ydbDriver, err := ydb.Unwrap(db.DB)
	if err != nil {
		return xerrors.WithStackTrace(err)
	}

	d.db = ydbDriver

	//connector, err := ydb.Connector(ydbDriver,
	//	ydb.WithQueryService(true),
	//	ydb.WithFakeTx(ydb.QueryExecuteQueryMode),
	//)
	//if err != nil {
	//	_ = ydbDriver.Close(context.Background())
	//
	//	return xerrors.WithStackTrace(err)
	//}
	//
	//db.DB = sql.OpenDB(connector)
	//
	//return d.Base.Init(core.FromDB(db.DB), d, uri, driverName, dataSource)

	return d.Base.Init(db, d, uri, driverName, dataSource)
}

func (d *YDB) GetIndexes(tableName string) (indexes map[string]*core.Index, err error) {
	end := d.log.Start("GetIndexes", tableName)
	defer func() { end.End(indexes, err) }()
	description, err := d.describeTable(context.Background(), tableName)
	if err != nil {
		return nil, xerrors.WithStackTrace(err)
	}

	return toCoreIndexes(description.Indexes), nil
}

func (d *YDB) IndexCheckSql(tableName, idxName string) (sql string, args []any) {
	end := d.log.Start("IndexCheckSql", tableName, idxName)
	defer func() { end.End(sql, args) }()

	return d.IndexCheckSQL(tableName, idxName)
}

func (d *YDB) SupportInsertMany() (ok bool) {
	end := d.log.Start("SupportInsertMany")
	defer func() { end.End(ok) }()
	return true // TODO:
}

func (d *YDB) SupportCharset() (ok bool) {
	end := d.log.Start("SupportCharset")
	defer func() { end.End(ok) }()
	return false
}

func (d *YDB) SupportEngine() (ok bool) {
	end := d.log.Start("SupportEngine")
	defer func() { end.End(ok) }()
	return false
}

func (d *YDB) SetParams(tableParams map[string]string) {
	end := d.log.Start("SetParams", tableParams)
	defer func() { end.End() }()
	d.tableParams = tableParams
}

func (d *YDB) absPath(tableName string) (res string) {
	end := d.log.Start("absPath", tableName)
	defer func() { end.End(res) }()
	return d.db.Name() + "/" + tableName
}

func (d *YDB) describeTable(
	ctx context.Context,
	tableName string,
) (description *options.Description, err error) {
	end := d.log.Start("describeTable", ctx, tableName)
	defer func() { end.End(description, err) }()
	description, err = d.db.Table().DescribeTable(ctx, d.absPath(tableName))
	if err != nil {
		return nil, xerrors.WithStackTrace(err)
	}

	return description, nil
}

func (d *YDB) IsTableExist(
	ctx context.Context,
	tableName string,
) (exists bool, err error) {
	end := d.log.Start("IsTableExist", ctx, tableName)
	defer func() { end.End(exists, err) }()
	exists, err = sugar.IsTableExists(ctx, d.db.Scheme(), d.absPath(tableName))
	if err != nil {
		return false, xerrors.WithStackTrace(err)
	}

	return exists, nil
}

func (d *YDB) TableCheckSql(tableName string) (sql string, args []any) {
	end := d.log.Start("TableCheckSql", tableName)
	defer func() { end.End(sql, args) }()
	return "SELECT Path FROM `.sys/partition_stats` WHERE Path LIKE '%/' || ?", []any{tableName}
}

func (d *YDB) AutoIncrStr() (res string) {
	end := d.log.Start("AutoIncrStr")
	defer func() { end.End(res) }()
	return ""
}

func (d *YDB) IsReserved(name string) (ok bool) {
	end := d.log.Start("IsReserved", name)
	defer func() { end.End(ok) }()
	_, ok = ydbReservedWords[strings.ToUpper(name)]
	return ok
}

func (d *YDB) SqlType(column *core.Column) (res string) {
	end := d.log.Start("SqlType", column)
	defer func() { end.End(res) }()
	return toYQLDataType(column.SQLType.Name, column.IsAutoIncrement)
}

// https://pkg.go.dev/database/sql#ColumnType.DatabaseTypeName
func (d *YDB) ColumnTypeKind(t string) (kind int) {
	end := d.log.Start("ColumnTypeKind", t)
	defer func() { end.End(kind) }()
	switch t {
	case types.TypeInt8.Yql(), types.TypeInt16.Yql(), types.TypeInt32.Yql(), types.TypeInt64.Yql(),
		types.TypeUint8.Yql(), types.TypeUint16.Yql(), types.TypeUint32.Yql(), types.TypeUint64.Yql():
		return core.NUMERIC_TYPE
	case types.TypeText.Yql(), types.TypeUTF8.Yql():
		return core.TEXT_TYPE
	case types.TypeTimestamp.Yql(), types.TypeTimestamp64.Yql(), types.TypeDatetime.Yql(), types.TypeDatetime64.Yql():
		return core.TIME_TYPE
	default:
		return core.UNKNOW_TYPE
	}
}

// YDB does not support this operation
func (d *YDB) ModifyColumnSQL(tableName string, column *core.Column) (res string) {
	end := d.log.Start("ModifyColumnSQL", tableName, column)
	defer func() { end.End(res) }()
	return ""
}

func (d *YDB) DropIndexSql(tableName string, index *core.Index) (res string) {
	end := d.log.Start("DropIndexSql", tableName, index)
	defer func() { end.End(res) }()
	var buf strings.Builder
	buf.WriteString("ALTER TABLE ")
	buf.WriteString(d.Quote(index.Name))
	buf.WriteString(" DROP INDEX ")
	buf.WriteString(d.Quote(tableName))
	buf.WriteString(";")
	return buf.String()
}

func (d *YDB) IndexOnTable() (ok bool) {
	end := d.log.Start("IndexOnTable")
	defer func() { end.End(ok) }()
	return true // TODO:
}

func (d *YDB) IsColumnExist(
	tableName,
	columnName string,
) (_ bool, err error) {
	description, err := d.describeTable(context.Background(), tableName)
	if err != nil {
		return false, xerrors.WithStackTrace(err)
	}

	return len(xslices.Filter(description.Columns, func(col options.Column) bool {
		return col.Name == columnName
	})) > 0, nil
}

func isJsonType(v types.Type) bool {
	isOptional, t := types.IsOptional(v)
	if !isOptional {
		t = v
	}

	if types.Equal(t, types.TypeJSONDocument) {
		return true
	}

	if types.Equal(t, types.TypeJSON) {
		return true
	}

	return false
}

func defaultValue(col options.Column) string {
	if col.DefaultValue == nil {
		return ""
	}
	v := col.DefaultValue.Literal()
	if v == nil {
		return ""
	}

	return v.Value.Yql()
}

func isAutoIncrement(col options.Column) bool {
	if col.DefaultValue == nil {
		return false
	}

	return col.DefaultValue.Sequence() != nil
}

func columnIndexes(indexes map[string]*core.Index, col options.Column) map[string]int {
	indexTypes := make(map[string]int, len(indexes))
	for _, idx := range xslices.Filter(xslices.Keys(indexes), func(idx string) bool {
		return slices.Contains(indexes[idx].Cols, col.Name)
	}) {
		indexTypes[idx] = indexes[idx].Type
	}

	return indexTypes
}

func toCoreColumn(
	tableName string,
	primaryKeys []string,
	tableIndexes map[string]*core.Index,
	col options.Column,
) *core.Column {
	isOptional, t := types.IsOptional(col.Type)
	if !isOptional {
		t = col.Type
	}

	d := defaultValue(col)

	return &core.Column{
		Name:            col.Name,
		TableName:       tableName,
		FieldName:       col.Name,
		SQLType:         yqlToSQLType(t.Yql()),
		IsJSON:          isJsonType(t),
		Nullable:        isOptional,
		Default:         d,
		Indexes:         columnIndexes(tableIndexes, col),
		IsPrimaryKey:    slices.Contains(primaryKeys, col.Name),
		IsAutoIncrement: isAutoIncrement(col),
		DefaultIsEmpty:  d == "",
		Comment:         col.Family,
	}
}

func toCoreIndexes(indexes []options.IndexDescription) map[string]*core.Index {
	return xslices.Map(
		xslices.Transform(indexes, toCoreIndex),
		func(idx *core.Index) string {
			return idx.Name
		},
	)
}

func toCoreIndexType(t options.IndexType) int {
	if t == options.IndexTypeGlobalUnique {
		return core.UniqueType
	}

	return core.IndexType
}

func toCoreIndex(d options.IndexDescription) *core.Index {
	return &core.Index{
		IsRegular: false,
		Name:      d.Name,
		Type:      toCoreIndexType(d.Type),
		Cols:      d.IndexColumns,
	}
}

func toCoreTable(d *options.Description) *core.Table {
	table := core.NewEmptyTable()
	table.Name = d.Name
	table.PrimaryKeys = d.PrimaryKey
	table.Indexes = toCoreIndexes(d.Indexes)

	for _, col := range d.Columns {
		table.AddColumn(
			toCoreColumn(
				d.Name,
				d.PrimaryKey,
				table.Indexes,
				col,
			),
		)
	}

	return table
}

func (d *YDB) GetColumns(tableName string) (
	names []string,
	cols map[string]*core.Column,
	err error,
) {
	end := d.log.Start("GetColumns", tableName)
	defer func() { end.End(names, cols, err) }()
	description, err := d.describeTable(context.Background(), tableName)
	if err != nil {
		return nil, nil, xerrors.WithStackTrace(err)
	}

	columns := xslices.Transform(description.Columns, func(col options.Column) *core.Column {
		return toCoreColumn(tableName, description.PrimaryKey, toCoreIndexes(description.Indexes), col)
	})

	return xslices.Transform(
			columns,
			func(col *core.Column) string {
				return col.Name
			},
		), xslices.Map(
			columns,
			func(col *core.Column) string {
				return col.Name
			},
		), nil
}

func (d *YDB) GetTables() (tables []*core.Table, err error) {
	end := d.log.Start("GetTables")
	defer func() {
		end.End(xslices.Transform(tables, func(t *core.Table) string {
			return t.Name
		}), err)
	}()
	ctx := context.Background()

	dir, err := d.db.Scheme().ListDirectory(ctx, d.db.Name())
	if err != nil {
		return nil, xerrors.WithStackTrace(err)
	}

	tables = make([]*core.Table, 0, len(dir.Children))
	for _, tableName := range xslices.Transform(
		xslices.Filter(
			dir.Children,
			func(e scheme.Entry) bool {
				return e.IsColumnTable() || e.IsTable()
			},
		),
		func(e scheme.Entry) string {
			return e.Name
		},
	) {
		description, err := d.describeTable(ctx, tableName)
		if err != nil {
			return nil, xerrors.WithStackTrace(err)
		}

		tables = append(tables, toCoreTable(description))
	}

	return tables, nil
}

// CreateTableSQLForCore generates `CREATE TABLE` YQL for xorm/core. Not part of core.Dialect.
func (d *YDB) CreateTableSQLForCore(
	ctx context.Context,
	_ any,
	table *core.Table,
	tableName string,
) (sql string, ok bool, err error) {
	end := d.log.Start("CreateTableSQLForCore", ctx, table, tableName)
	defer func() { end.End(sql, ok, err) }()
	tableName = d.Quote(tableName)

	var buf strings.Builder
	buf.WriteString("CREATE TABLE ")
	buf.WriteString(tableName)
	buf.WriteString(" ( ")

	// 	build primary key
	if len(table.PrimaryKeys) == 0 {
		return "", false, xerrors.WithStackTrace(errors.New("table must have at least one primary key"))
	}
	pk := make([]string, len(table.PrimaryKeys))
	pkMap := make(map[string]bool)
	for i := 0; i < len(table.PrimaryKeys); i++ {
		pk[i] = d.Quote(table.PrimaryKeys[i])
		pkMap[pk[i]] = true
	}
	var primaryKeyBuf strings.Builder
	primaryKeyBuf.WriteString("PRIMARY KEY ( ")
	primaryKeyBuf.WriteString(strings.Join(pk, ", "))
	primaryKeyBuf.WriteString(" )")
	primaryKey := primaryKeyBuf.String()

	// build column
	columnsList := make([]string, 0, len(table.Columns()))
	for _, c := range table.Columns() {
		columnName := d.Quote(c.Name)
		dataType := d.SqlType(c)
		if _, isPk := pkMap[columnName]; isPk {
			columnsList = append(columnsList, columnName+" "+dataType+" NOT NULL")
		} else {
			columnsList = append(columnsList, columnName+" "+dataType)
		}
	}
	joinColumns := strings.Join(columnsList, ", ")

	buf.WriteString(strings.Join([]string{joinColumns, primaryKey}, ", "))
	buf.WriteString(" ) ")

	if len(d.tableParams) > 0 {
		params := make([]string, 0, len(d.tableParams))
		for param, value := range d.tableParams {
			if param == "" || value == "" {
				continue
			}
			params = append(params, param+" = "+value)
		}
		if len(params) > 0 {
			buf.WriteString("WITH ( ")
			buf.WriteString(strings.Join(params, ", "))
			buf.WriteString(" ) ")
		}
	}

	buf.WriteString("; ")

	return buf.String(), true, nil
}

func (d *YDB) DropTableSQL(tableName string) (sql string, ok bool) {
	end := d.log.Start("DropTableSQL", tableName)
	defer func() { end.End(sql, ok) }()
	tableName = d.Quote(tableName)

	var buf strings.Builder
	buf.WriteString("DROP TABLE ")
	buf.WriteString(tableName)
	buf.WriteString(";")
	return buf.String(), false
}

func (d *YDB) Filters() (res []core.Filter) {
	return filters.Filters
}

func (d *YDB) IsRetryable(err error) (ok bool) {
	end := d.log.Start("IsRetryable", err)
	defer func() { end.End(ok) }()
	return retry.Check(err).MustRetry(false)
}

func (d *YDB) RetryOnError(err error) (ok bool) {
	end := d.log.Start("RetryOnError", err)
	defer func() { end.End(ok) }()
	return d.IsRetryable(err)
}
