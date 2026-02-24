package ydb

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/grafana/grafana/pkg/util/xorm/ydb/filters"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/pkg/xerrors"
	"github.com/ydb-platform/ydb-go-sdk/v3/retry"
	yc "github.com/ydb-platform/ydb-go-yc-metadata"

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
	ydbReservedWords = map[string]bool{
		"ABORT":             true,
		"ACTION":            true,
		"ADD":               true,
		"AFTER":             true,
		"ALL":               true,
		"ALTER":             true,
		"ANALYZE":           true,
		"AND":               true,
		"ANSI":              true,
		"ANY":               true,
		"ARRAY":             true,
		"AS":                true,
		"ASC":               true,
		"ASSUME":            true,
		"ASYNC":             true,
		"ATTACH":            true,
		"AUTOINCREMENT":     true,
		"AUTOMAP":           true,
		"BEFORE":            true,
		"BEGIN":             true,
		"BERNOULLI":         true,
		"BETWEEN":           true,
		"BITCAST":           true,
		"BY":                true,
		"CALLABLE":          true,
		"CASCADE":           true,
		"CASE":              true,
		"CAST":              true,
		"CHANGEFEED":        true,
		"CHECK":             true,
		"COLLATE":           true,
		"COLUMN":            true,
		"COLUMNS":           true,
		"COMMIT":            true,
		"COMPACT":           true,
		"CONDITIONAL":       true,
		"CONFLICT":          true,
		"CONSTRAINT":        true,
		"COVER":             true,
		"CREATE":            true,
		"CROSS":             true,
		"CUBE":              true,
		"CURRENT":           true,
		"CURRENT_TIME":      true,
		"CURRENT_DATE":      true,
		"CURRENT_TIMESTAMP": true,
		"DATABASE":          true,
		"DECIMAL":           true,
		"DECLARE":           true,
		"DEFAULT":           true,
		"DEFERRABLE":        true,
		"DEFERRED":          true,
		"DEFINE":            true,
		"DELETE":            true,
		"DESC":              true,
		"DETACH":            true,
		"DICT":              true,
		"DISABLE":           true,
		"DISCARD":           true,
		"DISTINCT":          true,
		"DO":                true,
		"DROP":              true,
		"EACH":              true,
		"ELSE":              true,
		"ERROR":             true,
		"EMPTY":             true,
		"EMPTY_ACTION":      true,
		"ENCRYPTED":         true,
		"END":               true,
		"ENUM":              true,
		"ERASE":             true,
		"ESCAPE":            true,
		"EVALUATE":          true,
		"EXCEPT":            true,
		"EXCLUDE":           true,
		"EXCLUSIVE":         true,
		"EXCLUSION":         true,
		"EXISTS":            true,
		"EXPLAIN":           true,
		"EXPORT":            true,
		"EXTERNAL":          true,
		"FAIL":              true,
		"FAMILY":            true,
		"FILTER":            true,
		"FLATTEN":           true,
		"FLOW":              true,
		"FOLLOWING":         true,
		"FOR":               true,
		"FOREIGN":           true,
		"FROM":              true,
		"FULL":              true,
		"FUNCTION":          true,
		"GLOB":              true,
		"GLOBAL":            true,
		"GROUP":             true,
		"GROUPING":          true,
		"GROUPS":            true,
		"HASH":              true,
		"HAVING":            true,
		"HOP":               true,
		"IF":                true,
		"IGNORE":            true,
		"ILIKE":             true,
		"IMMEDIATE":         true,
		"IMPORT":            true,
		"IN":                true,
		"INDEX":             true,
		"INDEXED":           true,
		"INHERITS":          true,
		"INITIALLY":         true,
		"INNER":             true,
		"INSERT":            true,
		"INSTEAD":           true,
		"INTERSECT":         true,
		"INTO":              true,
		"IS":                true,
		"ISNULL":            true,
		"JOIN":              true,
		"JSON_EXISTS":       true,
		"JSON_VALUE":        true,
		"JSON_QUERY":        true,
		"KEY":               true,
		"LEFT":              true,
		"LIKE":              true,
		"LIMIT":             true,
		"LIST":              true,
		"LOCAL":             true,
		"MATCH":             true,
		"NATURAL":           true,
		"NO":                true,
		"NOT":               true,
		"NOTNULL":           true,
		"NULL":              true,
		"NULLS":             true,
		"OBJECT":            true,
		"OF":                true,
		"OFFSET":            true,
		"ON":                true,
		"ONLY":              true,
		"OPTIONAL":          true,
		"OR":                true,
		"ORDER":             true,
		"OTHERS":            true,
		"OUTER":             true,
		"OVER":              true,
		"PARTITION":         true,
		"PASSING":           true,
		"PASSWORD":          true,
		"PLAN":              true,
		"PRAGMA":            true,
		"PRECEDING":         true,
		"PRESORT":           true,
		"PRIMARY":           true,
		"PROCESS":           true,
		"RAISE":             true,
		"RANGE":             true,
		"REDUCE":            true,
		"REFERENCES":        true,
		"REGEXP":            true,
		"REINDEX":           true,
		"RELEASE":           true,
		"RENAME":            true,
		"REPEATABLE":        true,
		"REPLACE":           true,
		"RESET":             true,
		"RESOURCE":          true,
		"RESPECT":           true,
		"RESTRICT":          true,
		"RESULT":            true,
		"RETURN":            true,
		"RETURNING":         true,
		"REVERT":            true,
		"RIGHT":             true,
		"RLIKE":             true,
		"ROLLBACK":          true,
		"ROLLUP":            true,
		"ROW":               true,
		"ROWS":              true,
		"SAMPLE":            true,
		"SAVEPOINT":         true,
		"SCHEMA":            true,
		"SELECT":            true,
		"SEMI":              true,
		"SET":               true,
		"SETS":              true,
		"STREAM":            true,
		"STRUCT":            true,
		"SUBQUERY":          true,
		"SYMBOLS":           true,
		"SYNC":              true,
		"SYSTEM":            true,
		"TABLE":             true,
		"TABLESAMPLE":       true,
		"TABLESTORE":        true,
		"TAGGED":            true,
		"TEMP":              true,
		"TEMPORARY":         true,
		"THEN":              true,
		"TIES":              true,
		"TO":                true,
		"TRANSACTION":       true,
		"TRIGGER":           true,
		"TUPLE":             true,
		"UNBOUNDED":         true,
		"UNCONDITIONAL":     true,
		"UNION":             true,
		"UNIQUE":            true,
		"UNKNOWN":           true,
		"UPDATE":            true,
		"UPSERT":            true,
		"USE":               true,
		"USER":              true,
		"USING":             true,
		"VACUUM":            true,
		"VALUES":            true,
		"VARIANT":           true,
		"VIEW":              true,
		"VIRTUAL":           true,
		"WHEN":              true,
		"WHERE":             true,
		"WINDOW":            true,
		"WITH":              true,
		"WITHOUT":           true,
		"WRAPPER":           true,
		"XOR":               true,
		"TRUE":              true,
		"FALSE":             true,
	}
)

const (
	// numeric types
	yql_Bool = "BOOL"

	yql_Int64 = "INT64"

	yql_Double = "DOUBLE"

	// serial types
	yql_Serial    = "SERIAL"
	yql_BigSerial = "BIGSERIAL"

	// string types
	yql_Bytes        = "BYTES"
	yql_Text         = "STRING"
	yql_JsonDocument = "JSONDOCUMENT"

	// Data and Time
	yql_DateTime  = "Datetime64"
	yql_Timestamp = "TIMESTAMP64"
)

func toYQLDataType(t string, isAutoIncrement bool) string {
	switch t {
	case core.Bool, core.Boolean:
		return yql_Bool
	case core.TinyInt:
		return yql_Int64
	case core.Int, core.Integer:
		if isAutoIncrement {
			return yql_Serial
		}
		return yql_Int64
	case core.SmallInt, core.MediumInt, core.BigInt:
		if isAutoIncrement {
			return yql_BigSerial
		}
		return yql_Int64
	case core.Float:
		return yql_Double
	case core.Double:
		return yql_Double
	case core.Blob, core.LongBlob, core.MediumBlob, core.TinyBlob, core.VarBinary, core.Binary:
		return yql_Bytes
	case core.Json:
		return yql_JsonDocument
	case core.Varchar, core.NVarchar, core.Char, core.NChar,
		core.MediumText, core.LongText, core.Text, core.NText, core.TinyText:
		return yql_Text
	case core.TimeStamp, core.Time, core.Date, core.DateTime:
		return yql_Timestamp
	case core.Serial:
		return yql_Serial
	case core.BigSerial:
		return yql_BigSerial
	default:
		return t
	}
}

func yqlToSQLType(yqlType string) core.SQLType {
	switch yqlType {
	case yql_Bool:
		return core.SQLType{Name: core.Bool, DefaultLength: 0, DefaultLength2: 0}
	case yql_Int64:
		return core.SQLType{Name: core.BigInt, DefaultLength: 0, DefaultLength2: 0}
	case yql_Double:
		return core.SQLType{Name: core.Double, DefaultLength: 0, DefaultLength2: 0}
	case yql_Bytes:
		return core.SQLType{Name: core.Blob, DefaultLength: 0, DefaultLength2: 0}
	case yql_JsonDocument:
		return core.SQLType{Name: core.Json, DefaultLength: 0, DefaultLength2: 0}
	case yql_Text:
		return core.SQLType{Name: core.Varchar, DefaultLength: 255, DefaultLength2: 0}
	case yql_Timestamp:
		return core.SQLType{Name: core.TimeStamp, DefaultLength: 0, DefaultLength2: 0}
	default:
		return core.SQLType{Name: yqlType}
	}
}

func removeOptional(s string) string {
	if s = strings.ToUpper(s); strings.HasPrefix(s, "OPTIONAL") {
		s = strings.TrimPrefix(s, "OPTIONAL<")
		s = strings.TrimSuffix(s, ">")
	}
	return s
}

var (
	_ core.Dialect                = (*Dialect)(nil)
	_ core.DialectWithReturningID = (*Dialect)(nil)
)

type Dialect struct {
	core.Base

	tableParams map[string]string // TODO: maybe remove
}

func (d *Dialect) WithReturningID() {}

func (d *Dialect) Init(db *core.DB, uri *core.Uri, drivername, dataSource string) error {
	ydbDriver, err := ydb.Unwrap(db.DB)
	if err != nil {
		return xerrors.WithStackTrace(err)
	}

	connector, err := ydb.Connector(ydbDriver,
		ydb.WithQueryService(true),
		ydb.WithFakeTx(ydb.QueryExecuteQueryMode),
	)
	if err != nil {
		_ = ydbDriver.Close(context.Background())

		return xerrors.WithStackTrace(err)
	}

	db.DB = sql.OpenDB(connector)

	return d.Base.Init(core.FromDB(db.DB), d, uri, drivername, dataSource)
}

func (d *Dialect) GetIndexes(tableName string) (map[string]*core.Index, error) {
	// YDB index introspection not implemented; return empty so xorm schema operations don't fail.
	return map[string]*core.Index{}, nil
}

func (d *Dialect) IndexCheckSql(tableName, idxName string) (string, []interface{}) {
	return "SELECT Path FROM `.sys/partition_stats` where Path LIKE '%/'" +
		" || ? || '/' || ? || '/indexImplTable'", []any{tableName, idxName}
}

func (d *Dialect) SupportInsertMany() bool {
	return true // TODO:
}

func (d *Dialect) SupportCharset() bool {
	return false
}

func (d *Dialect) SupportEngine() bool {
	return false
}

func (d *Dialect) WithConn(ctx context.Context, f func(context.Context, *sql.Conn) error) error {
	cc, err := d.DB().Conn(ctx)
	if err != nil {
		return xerrors.WithStackTrace(err)
	}
	defer func() {
		_ = cc.Close()
	}()

	return f(ctx, cc)
}

func (d *Dialect) WithConnRaw(ctx context.Context, f func(d interface{}) error) error {
	if err := d.WithConn(ctx, func(ctx context.Context, cc *sql.Conn) error {
		return cc.Raw(f)
	}); err != nil {
		return xerrors.WithStackTrace(err)
	}

	return nil
}

func (d *Dialect) SetParams(tableParams map[string]string) {
	d.tableParams = tableParams
}

func (d *Dialect) IsTableExist(
	ctx context.Context,
	tableName string,
) (_ bool, err error) {
	var exists bool
	err = d.WithConnRaw(ctx, func(dc interface{}) error {
		q, ok := dc.(interface {
			IsTableExists(context.Context, string) (bool, error)
		})
		if !ok {
			return xerrors.WithStackTrace(errors.New("driver hasn't method IsTableExists()"))
		}

		exists, err = q.IsTableExists(ctx, tableName)
		if err != nil {
			return xerrors.WithStackTrace(err)
		}

		return nil
	})
	if err != nil {
		return false, xerrors.WithStackTrace(err)
	}
	return exists, nil
}

func (d *Dialect) TableCheckSql(tableName string) (string, []any) {
	return "SELECT Path FROM `.sys/partition_stats` WHERE Path LIKE '%/' || ?", []any{tableName}
}

func (d *Dialect) AutoIncrStr() string {
	return ""
}

func (d *Dialect) IsReserved(name string) bool {
	_, ok := ydbReservedWords[strings.ToUpper(name)]
	return ok
}

func (d *Dialect) SqlType(column *core.Column) string {
	return toYQLDataType(column.SQLType.Name, column.IsAutoIncrement)
}

// sqlTypeForCreateTable returns the YQL type for a column when creating a table.
// migration_log.timestamp must be DATETIME so INSERT/RETURNING does not fail with "Timestamp to Optional<Date>".
func (d *Dialect) sqlTypeForCreateTable(tableName string, col *core.Column) string {
	unquoted := strings.Trim(tableName, "`\"")
	if strings.HasSuffix(unquoted, "migration_log") && col.Name == "timestamp" {
		return yql_DateTime
	}
	return d.SqlType(col)
}

// CreateTableSql implements the Dialect interface; used by Sync2 (e.g. migration_log).
// We override so migration_log.timestamp gets DATETIME (avoids "Timestamp to Optional<Date>" on INSERT).
func (d *Dialect) CreateTableSql(table *core.Table, tableName, storeEngine, charset string) string {
	if tableName == "" {
		tableName = table.Name
	}
	parts := make([]string, 0, len(table.ColumnsSeq())+2)
	for _, colName := range table.ColumnsSeq() {
		col := table.GetColumn(colName)
		dataType := d.sqlTypeForCreateTable(tableName, col)
		part := d.Quote(col.Name) + " " + dataType + " "
		if col.IsPrimaryKey && len(table.PrimaryKeys) == 1 {
			part += "PRIMARY KEY " + d.AutoIncrStr() + " "
		}
		if col.Default != "" {
			part += "DEFAULT " + col.Default + " "
		}
		if d.ShowCreateNull() {
			if col.Nullable {
				part += "NULL "
			} else {
				part += "NOT NULL "
			}
		}
		parts = append(parts, strings.TrimSpace(part))
	}
	if len(table.PrimaryKeys) > 1 {
		quoted := make([]string, len(table.PrimaryKeys))
		for i, pk := range table.PrimaryKeys {
			quoted[i] = d.Quote(pk)
		}
		parts = append(parts, "PRIMARY KEY ( "+strings.Join(quoted, ", ")+" )")
	}
	return "CREATE TABLE IF NOT EXISTS " + d.Quote(tableName) + " (" + strings.Join(parts, ", ") + ")"
}

// https://pkg.go.dev/database/sql#ColumnType.DatabaseTypeName
func (d *Dialect) ColumnTypeKind(t string) int {
	switch t {
	case "INT8", "INT16", "INT32", "INT64", "UINT8", "UINT16", "UINT32", "UINT64":
		return core.NUMERIC_TYPE
	case "UTF8":
		return core.TEXT_TYPE
	case "TIMESTAMP":
		return core.TIME_TYPE
	default:
		return core.UNKNOW_TYPE
	}
}

func (d *Dialect) Quote(name string) string {
	return "`" + name + "`"
}

func (d *Dialect) AddColumnSQL(tableName string, col *core.Column) string {
	tableName = d.Quote(tableName)
	columnName := d.Quote(col.Name)
	dataType := d.SqlType(col)

	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;", tableName, columnName, dataType))

	return buf.String()
}

// YDB does not support this operation
func (d *Dialect) ModifyColumnSQL(tableName string, column *core.Column) string {
	return ""
}

func (d *Dialect) DropIndexSql(tableName string, index *core.Index) string {
	tableName = d.Quote(tableName)
	indexName := d.Quote(index.Name)

	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("ALTER TABLE %s DROP INDEX %s;", tableName, indexName))

	return buf.String()
}

func (d *Dialect) IndexOnTable() bool {
	return true // TODO:
}

// TODO:
func (d *Dialect) IsColumnExist(
	tableName,
	columnName string,
) (_ bool, err error) {
	var exists bool
	ctx := context.TODO()
	err = d.WithConnRaw(ctx, func(dc interface{}) error {
		q, ok := dc.(interface {
			IsColumnExists(context.Context, string, string) (bool, error)
		})
		if !ok {
			return xerrors.WithStackTrace(errors.New("conn hasn't method IsColumnExists()"))
		}
		exists, err = q.IsColumnExists(ctx, tableName, columnName)
		if err != nil {
			return xerrors.WithStackTrace(err)
		}
		return nil
	})
	if err != nil {
		return false, xerrors.WithStackTrace(err)
	}
	return exists, nil
}

func (d *Dialect) GetColumns(tableName string) (
	_ []string,
	_ map[string]*core.Column,
	err error,
) {
	ctx := context.TODO()

	colNames := make([]string, 0)
	colMaps := make(map[string]*core.Column)

	err = d.WithConnRaw(ctx, func(dc interface{}) error {
		q, ok := dc.(interface {
			GetColumns(context.Context, string) ([]string, error)
			GetColumnType(context.Context, string, string) (string, error)
			IsPrimaryKey(context.Context, string, string) (bool, error)
		})
		if !ok {
			return xerrors.WithStackTrace(errors.New("driver does not support method [GetColumns]"))
		}

		colNames, err = q.GetColumns(ctx, tableName)
		if err != nil {
			return xerrors.WithStackTrace(err)
		}

		for _, colName := range colNames {
			dataType, err := q.GetColumnType(ctx, tableName, colName)
			if err != nil {
				return xerrors.WithStackTrace(err)
			}

			dataType = removeOptional(dataType)
			isPK, err := q.IsPrimaryKey(ctx, tableName, colName)
			if err != nil {
				return xerrors.WithStackTrace(err)
			}

			col := &core.Column{
				Name:         colName,
				TableName:    tableName,
				SQLType:      yqlToSQLType(dataType),
				IsPrimaryKey: isPK,
				Nullable:     !isPK,
				Indexes:      make(map[string]int),
			}
			if dataType == "SERIAL" {
				col.IsAutoIncrement = true
			}
			colMaps[colName] = col
		}
		return nil
	})
	if err != nil {
		return nil, nil, xerrors.WithStackTrace(err)
	}

	return colNames, colMaps, nil
}

func (d *Dialect) GetTables() (_ []*core.Table, err error) {
	tables := make([]*core.Table, 0)
	ctx := context.TODO()
	err = d.WithConnRaw(ctx, func(dc interface{}) error {
		q, ok := dc.(interface {
			GetTables(context.Context, string, bool, bool) ([]string, error)
		})
		if !ok {
			return xerrors.WithStackTrace(errors.New("driver does not support method [GetTables]"))
		}
		tableNames, err := q.GetTables(ctx, ".", true, true)
		if err != nil {
			return xerrors.WithStackTrace(err)
		}
		for _, tableName := range tableNames {
			table := core.NewEmptyTable()
			table.Name = tableName
			tables = append(tables, table)
		}
		return nil
	})
	if err != nil {
		return nil, xerrors.WithStackTrace(err)
	}
	return tables, nil
}

// !datbeohbbh! CreateTableSQL generate `CREATE TABLE` YQL.
// Method does not generate YQL for creating index.
func (d *Dialect) CreateTableSQL(
	ctx context.Context,
	_ any,
	table *core.Table,
	tableName string,
) (string, bool, error) {
	tableName = d.Quote(tableName)

	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("CREATE TABLE %s ( ", tableName))

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
	primaryKey := fmt.Sprintf("PRIMARY KEY ( %s )", strings.Join(pk, ", "))

	// build column
	columnsList := []string{}
	unquotedTableName := strings.Trim(tableName, "`\"")
	for _, c := range table.Columns() {
		columnName := d.Quote(c.Name)
		dataType := d.SqlType(c)
		// YDB INSERT/RETURNING for migration_log expects timestamp column as Date; TIMESTAMP64 causes "Timestamp to Optional<Date>" conversion error
		if strings.HasSuffix(unquotedTableName, "migration_log") && c.Name == "timestamp" {
			dataType = yql_DateTime
		}

		if _, isPk := pkMap[columnName]; isPk {
			columnsList = append(columnsList, fmt.Sprintf("%s %s NOT NULL", columnName, dataType))
		} else {
			columnsList = append(columnsList, fmt.Sprintf("%s %s", columnName, dataType))
		}
	}
	joinColumns := strings.Join(columnsList, ", ")

	buf.WriteString(strings.Join([]string{joinColumns, primaryKey}, ", "))
	buf.WriteString(" ) ")

	if len(d.tableParams) > 0 {
		params := make([]string, 0)
		for param, value := range d.tableParams {
			if param == "" || value == "" {
				continue
			}
			params = append(params, fmt.Sprintf("%s = %s", param, value))
		}
		if len(params) > 0 {
			buf.WriteString(fmt.Sprintf("WITH ( %s ) ", strings.Join(params, ", ")))
		}
	}

	buf.WriteString("; ")

	return buf.String(), true, nil
}

func (d *Dialect) DropTableSQL(tableName string) (string, bool) {
	tableName = d.Quote(tableName)

	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("DROP TABLE %s;", tableName))

	return buf.String(), false
}

func (d *Dialect) Filters() []core.Filter {
	return filters.Filters
}

func (d *Dialect) IsRetryable(err error) bool {
	return retry.Check(err).MustRetry(false)
}

func (d *Dialect) RetryOnError(err error) bool {
	return d.IsRetryable(err)
}
