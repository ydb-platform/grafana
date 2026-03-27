package migrator

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/url"
	"strconv"
	"strings"

	"github.com/lib/pq"
	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/retry"

	"xorm.io/core"
	"xorm.io/xorm"
)

var _ DialectRecursiveCTE = (*YDBDialect)(nil)

type YDBDialect struct {
	BaseDialect
}

func (d *YDBDialect) RecursiveQueriesAreSupported() (bool, error) {
	return false, nil
}

func NewYDBDialect() Dialect {
	d := YDBDialect{}
	d.dialect = &d
	d.driverName = YDB
	return &d
}

func (d *YDBDialect) IndexCheckSQL(tableName, indexName string) (string, []any) {
	return "SELECT Path FROM `.sys/partition_stats` where Path LIKE '%/'" +
		" || $1 || '/' || $2 || '/indexImplTable'", []any{tableName, indexName}
}

func (d *YDBDialect) SupportEngine() bool {
	return false
}

func (d *YDBDialect) Quote(name string) string {
	return "`" + name + "`"
}

func (d *YDBDialect) Concat(strs ...string) string {
	return strings.Join(strs, " || ")
}

func (d *YDBDialect) LikeOperator(column string, wildcardBefore bool, pattern string, wildcardAfter bool) (string, string) {
	param := pattern
	if wildcardBefore {
		param = "%" + param
	}
	if wildcardAfter {
		param = param + "%"
	}
	return fmt.Sprintf("%s ILIKE ?", column), param
}

func (d *YDBDialect) AutoIncrStr() string {
	return ""
}

func (d *YDBDialect) BooleanValue(value bool) any {
	return value
}

func (d *YDBDialect) BooleanStr(value bool) string {
	return strconv.FormatBool(value)
}

func (d *YDBDialect) BatchSize() int {
	return 1000
}

func (d *YDBDialect) SQLType(c *Column) string {
	xormDialect := core.QueryDialect(YDB)
	column := &core.Column{
		SQLType: core.SQLType{
			Name:           c.Type,
			DefaultLength:  c.Length,
			DefaultLength2: c.Length2,
		},
		IsAutoIncrement: c.IsAutoIncrement,
	}

	return xormDialect.SqlType(column)
}

func (d *YDBDialect) AddColumnSQL(tableName string, col *Column) string {
	col.Nullable = true // Cannot add not null column without default value
	col.Default = ""    // Column addition with default value is not supported now

	return d.BaseDialect.AddColumnSQL(tableName, col)
}

func (d *YDBDialect) RenameColumn(table Table, column *Column, newName string) string {
	oldName := column.Name
	column.Name = newName
	sql := d.AddColumnSQL(table.Name, column) + ";"
	if !column.IsPrimaryKey {
		column.Name = oldName
		sql += d.DropColumn(table, column)
	}

	return sql
}

// TODO:
func (d *YDBDialect) ColumnCheckSQL(tableName, columnName string) (string, []any) {
	return "", nil
}

func (d *YDBDialect) DropColumn(table Table, column *Column) string {
	return fmt.Sprintf("alter table %s DROP COLUMN %s", d.dialect.Quote(table.Name), d.dialect.Quote(column.Name))
}

func (d *YDBDialect) DropIndexSQL(tableName string, index *Index) string {
	return fmt.Sprintf("alter table %s DROP INDEX %s", d.dialect.Quote(tableName), d.dialect.Quote(index.XName(tableName)))
}

func (d *YDBDialect) UpdateTableSQL(tableName string, columns []*Column) string {
	return ""
	statements := []string{}

	for _, col := range columns {
		statements = append(statements, "ALTER "+d.Quote(col.Name)+" TYPE "+d.SQLType(col))
	}

	return "ALTER TABLE " + d.Quote(tableName) + " " + strings.Join(statements, ", ") + ";"
}

func (d *YDBDialect) CleanDB(engine *xorm.Engine) error {
	sess := engine.NewSession()
	defer sess.Close()

	if _, err := sess.Exec("DROP SCHEMA public CASCADE;"); err != nil {
		return fmt.Errorf("%v: %w", "failed to drop schema public", err)
	}

	if _, err := sess.Exec("CREATE SCHEMA public;"); err != nil {
		return fmt.Errorf("%v: %w", "failed to create schema public", err)
	}

	return nil
}

func (d *YDBDialect) Default(col *Column) string {
	if col.Type == DB_Bool {
		// Ensure that all dialects support the same literals in the same way.
		bl, err := strconv.ParseBool(col.Default)
		if err != nil {
			panic(fmt.Errorf("failed to create default value for column '%s': invalid boolean default value '%s'", col.Name, col.Default))
		}
		return d.dialect.BooleanStr(bl)
	}

	if col.Type == DB_NVarchar {
		return `"` + col.Default + `"`
	}

	return col.Default
}

func (d *YDBDialect) ColStringNoPk(col *Column) string {
	sql := d.dialect.Quote(col.Name) + " "

	sql += d.dialect.SQLType(col) + " NULL " // TODO: remove always NULL when done with add not null columns

	if col.Default != "" {
		sql += "DEFAULT " + d.dialect.Default(col) + " "
	}

	return sql
}

func (d *YDBDialect) CreateTableSQL(table *Table) string {
	sql := "CREATE TABLE IF NOT EXISTS "
	sql += d.dialect.Quote(table.Name) + " (\n"

	pkList := table.PrimaryKeys

	for _, col := range table.Columns {
		if len(pkList) == 0 && !col.Nullable {
			pkList = []string{col.Name}
		}

		sql += col.StringNoPk(d.dialect)
		sql = strings.TrimSpace(sql)
		sql += "\n, "
	}

	quotedCols := []string{}
	for _, col := range pkList {
		quotedCols = append(quotedCols, d.dialect.Quote(col))
	}

	sql += "PRIMARY KEY ( " + strings.Join(quotedCols, ",") + " ), "

	sql = sql[:len(sql)-2] + ")"

	sql += ";"
	return sql
}

// TruncateDBTables truncates all the tables.
// A special case is the dashboard_acl table where we keep the default permissions.
func (d *YDBDialect) TruncateDBTables(engine *xorm.Engine) error {
	tables, err := engine.Dialect().GetTables()
	if err != nil {
		return err
	}
	sess := engine.NewSession()
	defer sess.Close()

	dbName, err := d.GetDBName(engine.DataSourceName())
	if err != nil {
		return err
	}

	for _, table := range tables {
		switch table.Name {
		case "":
			continue
		case "migration_log":
			continue
		case "dashboard_acl":
			// keep default dashboard permissions
			if _, err := sess.Exec(fmt.Sprintf("DELETE FROM %v WHERE dashboard_id != -1 AND org_id != -1;", d.Quote(table.Name))); err != nil {
				return fmt.Errorf("failed to truncate table %q: %w", table.Name, err)
			}
			if _, err := sess.Exec(fmt.Sprintf("ALTER SEQUENCE %v RESTART WITH 3;", d.Quote(fmt.Sprintf("%s/%v/_serial_column_id", dbName, table.Name)))); err != nil {
				return fmt.Errorf("failed to reset table %q: %w", table.Name, err)
			}
		default:
			err := retry.Do(context.Background(), engine.DB().DB, func(ctx context.Context, cc *sql.Conn) error {
				_, err := sess.Exec(fmt.Sprintf("DELETE FROM %v;", d.Quote(table.Name)))
				return err
			})
			if err != nil {
				if d.isUndefinedTable(err) {
					continue
				}
				return fmt.Errorf("failed to truncate table %q: %w", table.Name, err)
			}

			_, tableCols, err := engine.Dialect().GetColumns(table.Name)
			if err != nil {
				return err
			}

			for _, column := range tableCols {
				if column.IsAutoIncrement {
					sequenceName := fmt.Sprintf("%v/%v/_serial_column_%v", dbName, table.Name, column.Name)
					if _, err := sess.Exec(fmt.Sprintf("ALTER SEQUENCE %v RESTART;", d.Quote(sequenceName))); err != nil {
						return fmt.Errorf("failed to reset sequence %q: %w", sequenceName, err)
					}
				}
			}
		}
	}

	return nil
}

func (d *YDBDialect) isThisError(err error, errcode string) bool {
	var driverErr *pq.Error
	if errors.As(err, &driverErr) {
		if string(driverErr.Code) == errcode {
			return true
		}
	}

	return false
}

func (d *YDBDialect) ErrorMessage(err error) string {
	var driverErr *pq.Error
	if errors.As(err, &driverErr) {
		return driverErr.Message
	}
	return ""
}

func (d *YDBDialect) isUndefinedTable(err error) bool {
	return ydb.IsOperationErrorSchemeError(err)
}

func (d *YDBDialect) IsUniqueConstraintViolation(err error) bool {
	return d.isThisError(err, "23505")
}

func (d *YDBDialect) IsDeadlock(err error) bool {
	return d.isThisError(err, "40P01")
}

func (d *YDBDialect) CreateIndexSQL(tableName string, index *Index) string {
	indexName := d.Quote(index.XName(tableName))
	tableName = d.Quote(tableName)

	colsIndex := make([]string, len(index.Cols))
	for i := 0; i < len(index.Cols); i++ {
		colsIndex[i] = d.Quote(index.Cols[i])
	}

	indexOn := strings.Join(colsIndex, ",")

	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("ALTER TABLE %s ADD INDEX %s GLOBAL ON ( %s );", tableName, indexName, indexOn))

	return buf.String()
}

// UpsertSQL returns the upsert sql statement for PostgreSQL dialect
func (d *YDBDialect) UpsertSQL(tableName string, keyCols, updateCols []string) string {
	str, _ := d.UpsertMultipleSQL(tableName, keyCols, updateCols, 1)
	return str
}

// UpsertMultipleSQL returns the upsert sql statement for PostgreSQL dialect
func (d *YDBDialect) UpsertMultipleSQL(tableName string, keyCols, updateCols []string, count int) (string, error) {
	if count < 1 {
		return "", fmt.Errorf("upsert statement must have count >= 1. Got %v", count)
	}
	columnsStr := strings.Builder{}
	const separator = ", "
	separatorVar := separator
	for i, c := range updateCols {
		if i == len(updateCols)-1 {
			separatorVar = ""
		}
		columnsStr.WriteString(fmt.Sprintf("%s%s", d.Quote(c), separatorVar))
	}

	valuesStr := strings.Builder{}
	separatorVar = separator
	nextPlaceHolder := 1

	for i := 0; i < count; i++ {
		if i == count-1 {
			separatorVar = ""
		}

		colPlaceHoldersStr := strings.Builder{}
		placeHolderSep := separator
		for j := 1; j <= len(updateCols); j++ {
			if j == len(updateCols) {
				placeHolderSep = ""
			}
			placeHolder := fmt.Sprintf("$%v%s", nextPlaceHolder, placeHolderSep)
			nextPlaceHolder++
			colPlaceHoldersStr.WriteString(placeHolder)
		}
		colPlaceHolders := colPlaceHoldersStr.String()

		valuesStr.WriteString(fmt.Sprintf("(%s)%s", colPlaceHolders, separatorVar))
	}

	s := fmt.Sprintf(`UPSERT INTO %s (%s) VALUES %s;`,
		tableName,
		columnsStr.String(),
		valuesStr.String(),
	)

	return s, nil
}

func (d *YDBDialect) GetDBName(dsn string) (string, error) {
	uri, err := url.Parse(dsn)
	if err != nil {
		return "", fmt.Errorf("failed on parse data source %v", dsn)
	}

	return uri.Path, nil
}

// OrderBy returns ORDER BY expression for YDB. The search subquery only selects id,
// so "dashboard.title" is not in the subquery result and causes "Member not found".
// We use "title" and the builder must add "dashboard.title AS title" to the subquery
// SELECT when ordering by title (see searchstore/builder.go).
func (d *YDBDialect) OrderBy(order string) string {
	order = strings.ReplaceAll(order, "dashboard.title", "title")
	return order
}
