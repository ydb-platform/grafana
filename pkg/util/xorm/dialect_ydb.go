package xorm

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/retry"
	yc "github.com/ydb-platform/ydb-go-yc-metadata"

	"github.com/grafana/grafana/pkg/util/xorm/core"
)

// type ydbDriver struct {
// }

// func (p *ydbDriver) Parse(driverName, dataSourceName string) (*core.Uri, error) {
// 	// if strings.Contains(dataSourceName, "?") {
// 	// 	dataSourceName = dataSourceName[:strings.Index(dataSourceName, "?")]
// 	// }

// 	return &core.Uri{DbType: core.YDB, DbName: dataSourceName}, nil
// }

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

	// ydbQuoter = core.Quoter{
	// 	Prefix:     '`',
	// 	Suffix:     '`',
	// 	IsReserved: core.AlwaysReserve,
	// }
)

const (
	// numeric types
	yql_Bool = "BOOL"

	yql_Int8  = "INT8"
	yql_Int16 = "INT16"
	yql_Int32 = "INT32"
	yql_Int64 = "INT64"

	yql_Uint8  = "UINT8"
	yql_Uint16 = "UINT16"
	yql_Uint32 = "UINT32"
	yql_Uint64 = "UINT64"

	yql_Float   = "FLOAT"
	yql_Double  = "DOUBLE"
	yql_Decimal = "DECIMAL"

	// serial types
	yql_Serial    = "SERIAL"
	yql_BigSerial = "BIGSERIAL"

	// string types
	yql_String       = "STRING"
	yql_Utf8         = "UTF8"
	yql_Json         = "JSON"
	yql_JsonDocument = "JSONDOCUMENT"
	yql_Yson         = "YSON"

	// Data and Time
	yql_Date      = "DATE"
	yql_DateTime  = "DATETIME"
	yql_Timestamp = "TIMESTAMP64"
	yql_Interval  = "INTERVAL"

	// Containers
	yql_List = "LIST"
)

func toYQLDataType(t string, isAutoIncrement bool) string {
	switch t {
	case core.Bool, core.Boolean:
		return yql_Bool
	case core.TinyInt:
		return yql_Int8
	case core.Int, core.Integer:
		if isAutoIncrement {
			return yql_Serial
		}
		return yql_Int32
	case core.SmallInt, core.MediumInt, core.BigInt:
		if isAutoIncrement {
			return yql_BigSerial
		}
		return yql_Int64
	case core.Float:
		return yql_Float
	case core.Double:
		return yql_Double
	case core.Blob, core.LongBlob, core.MediumBlob, core.TinyBlob, core.VarBinary, core.Binary:
		return yql_String
	case core.Json:
		return yql_Json
	case core.Varchar, core.NVarchar, core.Char, core.NChar,
		core.MediumText, core.LongText, core.Text, core.NText, core.TinyText:
		return yql_Utf8
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
	case yql_Int8:
		return core.SQLType{Name: core.TinyInt, DefaultLength: 0, DefaultLength2: 0}
	case yql_Int16:
		return core.SQLType{Name: core.SmallInt, DefaultLength: 0, DefaultLength2: 0}
	case yql_Int32:
		return core.SQLType{Name: core.Integer, DefaultLength: 0, DefaultLength2: 0}
	case yql_Int64:
		return core.SQLType{Name: core.BigInt, DefaultLength: 0, DefaultLength2: 0}
	case yql_Float:
		return core.SQLType{Name: core.Float, DefaultLength: 0, DefaultLength2: 0}
	case yql_Double:
		return core.SQLType{Name: core.Double, DefaultLength: 0, DefaultLength2: 0}
	case yql_String:
		return core.SQLType{Name: core.Blob, DefaultLength: 0, DefaultLength2: 0}
	case yql_Json:
		return core.SQLType{Name: core.Json, DefaultLength: 0, DefaultLength2: 0}
	case yql_Utf8:
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

type ydbDialect struct {
	core.Base

	tableParams map[string]string // TODO: maybe remove
}

// ydbConnectorWrapper wraps driver.Connector to intercept connection creation
type ydbConnectorWrapper struct {
	connector driver.Connector
}

// ydbConnWrapper wraps driver.Conn to intercept statement preparation
type ydbConnWrapper struct {
	conn driver.Conn
}

// ydbStmtWrapper wraps driver.Stmt to convert time.Duration arguments to int64
type ydbStmtWrapper struct {
	stmt  driver.Stmt
	query string
}

const ydbCostBasedOptimizationPragma = `PRAGMA ydb.CostBasedOptimization = "on";` + "\n"

func prependYdbPragma(query string) string {
	trimmed := strings.TrimSpace(query)
	if trimmed == "" {
		return query
	}
	if strings.HasPrefix(trimmed, strings.TrimSpace(ydbCostBasedOptimizationPragma)) {
		return query
	}
	return ydbCostBasedOptimizationPragma + query
}

// migration todo:
//
// ALTER TABLE `cache_data` ADD COLUMN created_at_new Int64;
// UPDATE `cache_data` SET created_at_new = CAST(created_at AS Int64);
// ALTER TABLE `cache_data` DROP COLUMN created_at;
// ALTER TABLE `cache_data` ADD COLUMN created_at Int64;
// UPDATE `cache_data` SET created_at = created_at_new;
// ALTER TABLE `cache_data` DROP COLUMN created_at_new;
//
// ALTER TABLE `cache_data` ADD COLUMN expires_new Int64;
// UPDATE `cache_data` SET expires_new = CAST(expires AS Int64);
// ALTER TABLE `cache_data` DROP COLUMN expires;
// ALTER TABLE `cache_data` ADD COLUMN expires Int64;
// UPDATE `cache_data` SET expires = expires_new;
// ALTER TABLE `cache_data` DROP COLUMN expires_new;
//
// ALTER TABLE `user` ADD COLUMN version_new Int64;
// UPDATE `user` SET version_new = CAST(version AS Int64);
// ALTER TABLE `user` DROP COLUMN version;
// ALTER TABLE `user` ADD COLUMN version Int64;
// UPDATE `user` SET version = version_new;
// ALTER TABLE `user` DROP COLUMN version_new;

// convertArgs converts time.Duration arguments to int64
func convertArgs(args []driver.NamedValue) []driver.NamedValue {
	converted := make([]driver.NamedValue, len(args))
	for i, arg := range args {
		converted[i] = arg
		if duration, ok := arg.Value.(time.Duration); ok {
			converted[i].Value = int64(duration)
		}
	}
	return converted
}

// parseInsertColumnOrdinal returns the 1-based ordinal of a column in the column list
// of an INSERT/REPLACE query, or 0 if the query is not recognized or the column is not found.
// Allows converting arguments without binding to specific table or column names.
func parseInsertColumnOrdinal(query, columnName string) int {
	queryLower := strings.ToLower(strings.TrimSpace(query))
	insertIdx := strings.Index(queryLower, "insert into")
	if insertIdx < 0 {
		insertIdx = strings.Index(queryLower, "replace into")
	}
	if insertIdx < 0 {
		return 0
	}
	// Find the first '(' after INSERT INTO / REPLACE INTO — that is the column list
	afterInsert := query[insertIdx:]
	open := strings.Index(afterInsert, "(")
	if open < 0 {
		return 0
	}
	// Find the matching closing parenthesis
	depth := 1
	for i := open + 1; i < len(afterInsert); i++ {
		switch afterInsert[i] {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				columnList := afterInsert[open+1 : i]
				columns := splitSQLColumnList(columnList)
				colLower := strings.ToLower(columnName)
				for idx, col := range columns {
					if strings.ToLower(col) == colLower {
						return idx + 1
					}
				}
				return 0
			}
		}
	}
	return 0
}

// splitSQLColumnList splits a column list like "`a`, `b`, `c`" or "a, b, c" into column names.
func splitSQLColumnList(list string) []string {
	var columns []string
	var current strings.Builder
	inQuote := false
	var quoteChar rune
	for _, r := range list {
		switch {
		case !inQuote && (r == '`' || r == '"'):
			inQuote = true
			quoteChar = r
		case inQuote && r == quoteChar:
			inQuote = false
		case !inQuote && (r == ',' || r == ' '):
			if r == ',' {
				if s := strings.TrimSpace(current.String()); s != "" {
					columns = append(columns, s)
				}
				current.Reset()
			}
		default:
			current.WriteRune(r)
		}
	}
	if s := strings.TrimSpace(current.String()); s != "" {
		columns = append(columns, s)
	}
	return columns
}

// parseUpdateColumnOrdinal returns the 1-based ordinal of the placeholder for columnName in the
// SET clause of an UPDATE query (e.g. UPDATE t SET a=?, b=?, c=? -> ordinal of "b" is 2), or 0 if not found.
func parseUpdateColumnOrdinal(query, columnName string) int {
	queryLower := strings.ToLower(strings.TrimSpace(query))
	setIdx := strings.Index(queryLower, " set ")
	if setIdx < 0 {
		return 0
	}
	// SET clause: from " SET " to " WHERE " or end of query
	afterSet := query[setIdx+5:]
	whereIdx := strings.Index(strings.ToLower(afterSet), " where ")
	var setClause string
	if whereIdx >= 0 {
		setClause = strings.TrimSpace(afterSet[:whereIdx])
	} else {
		setClause = strings.TrimSpace(afterSet)
	}
	// Split by comma to get "col=?" parts (no nested parens in typical SET)
	var columns []string
	var current strings.Builder
	depth := 0
	for _, r := range setClause {
		switch r {
		case '(':
			depth++
			current.WriteRune(r)
		case ')':
			depth--
			current.WriteRune(r)
		case ',':
			if depth == 0 {
				part := strings.TrimSpace(current.String())
				if part != "" {
					col := columnNameFromSetPart(part)
					if col != "" {
						columns = append(columns, col)
					}
				}
				current.Reset()
			} else {
				current.WriteRune(r)
			}
		default:
			current.WriteRune(r)
		}
	}
	if s := strings.TrimSpace(current.String()); s != "" {
		if col := columnNameFromSetPart(s); col != "" {
			columns = append(columns, col)
		}
	}
	colLower := strings.ToLower(columnName)
	for idx, col := range columns {
		if strings.ToLower(col) == colLower {
			return idx + 1
		}
	}
	return 0
}

// columnNameFromSetPart extracts the column name from "col=?" or "col = ?" or "`col`=?".
func columnNameFromSetPart(part string) string {
	eq := strings.Index(part, "=")
	if eq <= 0 {
		return ""
	}
	name := strings.TrimSpace(part[:eq])
	// Strip quotes/backticks
	name = strings.Trim(name, "`\"")
	return strings.TrimSpace(name)
}

// isInt64ToInt32ConversionError reports whether the error is an Int64 -> Int32 conversion error.
func isInt64ToInt32ConversionError(err error) bool {
	if err == nil {
		return false
	}
	errStr := err.Error()
	return strings.Contains(errStr, "Failed to convert") &&
		(strings.Contains(errStr, "Int64 to Int32") || strings.Contains(errStr, "Int64 to Optional<Int32>"))
}

// extractFieldNameFromError extracts the field name from a conversion error.
// Error format: "Failed to convert 'version': Int64 to Int32" or "Failed to convert 'version': Int64 to Optional<Int32>"
func extractFieldNameFromError(err error) string {
	if err == nil {
		return ""
	}
	errStr := err.Error()

	// Match pattern: 'Failed to convert '<field_name>': Int64 to [Optional<]Int32[>]'
	re := regexp.MustCompile(`Failed to convert '([^']+)': Int64 to (?:Optional<)?Int32(?:>)?`)
	matches := re.FindStringSubmatch(errStr)
	if len(matches) > 1 {
		return matches[1]
	}
	return ""
}

// convertInt64ToInt32Args converts the given field from int64 to int32.
// fieldName is taken from the YDB error; conversion is by parameter name or by position
// in the query (parsing INSERT/REPLACE), without binding to specific tables.
func convertInt64ToInt32Args(args []driver.NamedValue, fieldName, query string) []driver.NamedValue {
	converted := make([]driver.NamedValue, len(args))
	copy(converted, args)

	if fieldName == "" {
		return converted
	}

	fieldLower := strings.ToLower(fieldName)
	convertedByOrdinal := false

	for i := range converted {
		arg := &converted[i]
		argNameLower := strings.ToLower(arg.Name)
		if argNameLower == fieldLower ||
			strings.Contains(argNameLower, fieldLower) ||
			strings.HasSuffix(argNameLower, "."+fieldLower) {
			if v, ok := arg.Value.(int64); ok {
				arg.Value = int32(v)
				convertedByOrdinal = true
			}
		}
	}

	// If not found by parameter name — determine position from query text (INSERT or UPDATE)
	if !convertedByOrdinal && query != "" {
		ord := parseInsertColumnOrdinal(query, fieldName)
		if ord == 0 {
			ord = parseUpdateColumnOrdinal(query, fieldName)
		}
		if ord > 0 {
			// Prefer conversion by Ordinal (1-based) in case slice order differs from placeholder order
			foundByOrdinal := false
			for i := range converted {
				if converted[i].Ordinal == ord {
					if v, ok := converted[i].Value.(int64); ok {
						converted[i].Value = int32(v)
					}
					foundByOrdinal = true
					break
				}
			}
			if !foundByOrdinal && ord <= len(converted) {
				if v, ok := converted[ord-1].Value.(int64); ok {
					converted[ord-1].Value = int32(v)
				}
			}
		}
	}

	return converted
}

// Connect implements driver.Connector interface
func (w *ydbConnectorWrapper) Connect(ctx context.Context) (driver.Conn, error) {
	conn, err := w.connector.Connect(ctx)
	if err != nil {
		return nil, err
	}
	return &ydbConnWrapper{conn: conn}, nil
}

// Driver implements driver.Connector interface
func (w *ydbConnectorWrapper) Driver() driver.Driver {
	return w.connector.Driver()
}

// CheckNamedValue implements driver.NamedValueChecker interface
func (w *ydbConnWrapper) CheckNamedValue(nv *driver.NamedValue) error {
	// Convert time.Duration to int64
	if duration, ok := nv.Value.(time.Duration); ok {
		nv.Value = int64(duration)
		// if nv.Name == "created_at" {
		// 	nv.Value = int32(duration)
		// }
	}

	rv := reflect.ValueOf(nv.Value)
	if rv.Kind() == reflect.Int {
		nv.Value = rv.Int()
	}

	// Handle nil time values
	// if nv.Value == nil && nv.Ordinal == 10 {
	// 	var nilTime *time.Time
	// 	nv.Value = nilTime
	// 	return nil
	// }

	// Delegate to underlying connector if it implements NamedValueChecker
	if checker, ok := w.conn.(driver.NamedValueChecker); ok {
		return checker.CheckNamedValue(nv)
	}

	return nil
}

// Prepare implements driver.Conn interface
func (w *ydbConnWrapper) Prepare(query string) (driver.Stmt, error) {
	// query = prependYdbPragma(query)
	stmt, err := w.conn.Prepare(query)
	if err != nil {
		return nil, err
	}
	return &ydbStmtWrapper{stmt: stmt, query: query}, nil
}

// PrepareContext implements driver.ConnPrepareContext interface
func (w *ydbConnWrapper) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	if connCtx, ok := w.conn.(driver.ConnPrepareContext); ok {
		// query = prependYdbPragma(query)
		stmt, err := connCtx.PrepareContext(ctx, query)
		if err != nil {
			return nil, err
		}
		return &ydbStmtWrapper{stmt: stmt, query: query}, nil
	}
	return w.Prepare(query)
}

// Close implements driver.Conn interface
func (w *ydbConnWrapper) Close() error {
	return w.conn.Close()
}

// Begin implements driver.Conn interface
func (w *ydbConnWrapper) Begin() (driver.Tx, error) {
	return w.conn.Begin()
}

// BeginTx implements driver.ConnBeginTx interface
func (w *ydbConnWrapper) BeginTx(ctx context.Context, opts driver.TxOptions) (driver.Tx, error) {
	if connTx, ok := w.conn.(driver.ConnBeginTx); ok {
		return connTx.BeginTx(ctx, opts)
	}
	return w.conn.Begin()
}

// GetTables delegates to the underlying connection's GetTables method
func (w *ydbConnWrapper) GetTables(ctx context.Context, path string, recursive bool, excludeSysTables bool) ([]string, error) {
	if connWithGetTables, ok := w.conn.(interface {
		GetTables(context.Context, string, bool, bool) ([]string, error)
	}); ok {
		return connWithGetTables.GetTables(ctx, path, recursive, excludeSysTables)
	}
	return nil, fmt.Errorf("underlying connection does not support GetTables method")
}

// IsTableExists delegates to the underlying connection's IsTableExists method
func (w *ydbConnWrapper) IsTableExists(ctx context.Context, tableName string) (bool, error) {
	if connWithIsTableExists, ok := w.conn.(interface {
		IsTableExists(context.Context, string) (bool, error)
	}); ok {
		return connWithIsTableExists.IsTableExists(ctx, tableName)
	}
	return false, fmt.Errorf("underlying connection does not support IsTableExists method")
}

// GetColumns delegates to the underlying connection's GetColumns method
func (w *ydbConnWrapper) GetColumns(ctx context.Context, tableName string) ([]string, error) {
	if connWithGetColumns, ok := w.conn.(interface {
		GetColumns(context.Context, string) ([]string, error)
	}); ok {
		return connWithGetColumns.GetColumns(ctx, tableName)
	}

	return nil, fmt.Errorf("underlying connection does not support GetColumns method")
}

func getLastPartOfColumn(column string) string {
	for i := len(column) - 1; i >= 0; i-- {
		if column[i] == '.' {
			return column[i+1:]
		}
	}

	return column
}

// GetColumnType delegates to the underlying connection's GetColumnType method
func (w *ydbConnWrapper) GetColumnType(ctx context.Context, tableName, columnName string) (string, error) {
	if connWithGetColumnType, ok := w.conn.(interface {
		GetColumnType(context.Context, string, string) (string, error)
	}); ok {
		return connWithGetColumnType.GetColumnType(ctx, tableName, columnName)
	}
	return "", fmt.Errorf("underlying connection does not support GetColumnType method")
}

// IsPrimaryKey delegates to the underlying connection's IsPrimaryKey method
func (w *ydbConnWrapper) IsPrimaryKey(ctx context.Context, tableName, columnName string) (bool, error) {
	if connWithIsPrimaryKey, ok := w.conn.(interface {
		IsPrimaryKey(context.Context, string, string) (bool, error)
	}); ok {
		return connWithIsPrimaryKey.IsPrimaryKey(ctx, tableName, columnName)
	}
	return false, fmt.Errorf("underlying connection does not support IsPrimaryKey method")
}

// IsColumnExists delegates to the underlying connection's IsColumnExists method
func (w *ydbConnWrapper) IsColumnExists(ctx context.Context, tableName, columnName string) (bool, error) {
	if connWithIsColumnExists, ok := w.conn.(interface {
		IsColumnExists(context.Context, string, string) (bool, error)
	}); ok {
		return connWithIsColumnExists.IsColumnExists(ctx, tableName, columnName)
	}
	return false, fmt.Errorf("underlying connection does not support IsColumnExists method")
}

// GetIndexes delegates to the underlying connection's GetIndexes method
func (w *ydbConnWrapper) GetIndexes(ctx context.Context, tableName string) ([]string, error) {
	if connWithGetIndexes, ok := w.conn.(interface {
		GetIndexes(context.Context, string) ([]string, error)
	}); ok {
		return connWithGetIndexes.GetIndexes(ctx, tableName)
	}
	return nil, fmt.Errorf("underlying connection does not support GetIndexes method")
}

// GetIndexColumns delegates to the underlying connection's GetIndexColumns method
func (w *ydbConnWrapper) GetIndexColumns(ctx context.Context, tableName, indexName string) ([]string, error) {
	if connWithGetIndexColumns, ok := w.conn.(interface {
		GetIndexColumns(context.Context, string, string) ([]string, error)
	}); ok {
		return connWithGetIndexColumns.GetIndexColumns(ctx, tableName, indexName)
	}
	return nil, fmt.Errorf("underlying connection does not support GetIndexColumns method")
}

type rowsWrapper struct {
	driver.Rows
}

func (w *rowsWrapper) Columns() []string {
	columns := w.Rows.Columns()
	for i := range columns {
		columns[i] = getLastPartOfColumn(columns[i])
	}

	return columns
}

// QueryContext intercepts query execution and converts time.Duration to int64
func (w *ydbStmtWrapper) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	args = convertArgs(args)

	// Handle special case for LIMIT clause
	if strings.HasSuffix(w.query, "LIMIT $3;\n") {
		for i, arg := range args {
			if arg.Ordinal == 3 {
				if val, ok := arg.Value.(int64); ok {
					args[i].Value = uint64(val)
				}
			}
		}
	}

	// Execute with retry loop: one INSERT/UPDATE can have multiple Int32 columns (version, updated_by, etc.)
	if stmtCtx, ok := w.stmt.(driver.StmtQueryContext); ok {
		const maxInt32Retries = 20
		currentArgs := args
		for attempt := 0; attempt < maxInt32Retries; attempt++ {
			rows, err := stmtCtx.QueryContext(ctx, currentArgs)
			if err == nil {
				return &rowsWrapper{Rows: rows}, nil
			}
			if !isInt64ToInt32ConversionError(err) {
				return nil, err
			}
			fieldName := extractFieldNameFromError(err)
			currentArgs = convertInt64ToInt32Args(currentArgs, fieldName, w.query)
		}
		return nil, fmt.Errorf("Int64 to Int32 conversion retry limit (%d) exceeded", maxInt32Retries)
	}

	// Fallback to non-context query with retry loop
	currentArgs := args
	const maxInt32Retries = 20
	for attempt := 0; attempt < maxInt32Retries; attempt++ {
		values := make([]driver.Value, len(currentArgs))
		for i, arg := range currentArgs {
			values[i] = arg.Value
		}
		rows, err := w.stmt.Query(values)
		if err == nil {
			return rows, nil
		}
		if !isInt64ToInt32ConversionError(err) {
			return nil, err
		}
		fieldName := extractFieldNameFromError(err)
		currentArgs = convertInt64ToInt32Args(currentArgs, fieldName, w.query)
	}
	return nil, fmt.Errorf("Int64 to Int32 conversion retry limit (%d) exceeded", maxInt32Retries)
}

// ExecContext intercepts exec execution and converts time.Duration to int64
func (w *ydbStmtWrapper) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	args = convertArgs(args)

	// Execute with retry loop: one INSERT/UPDATE can have multiple Int32 columns
	if stmtCtx, ok := w.stmt.(driver.StmtExecContext); ok {
		const maxInt32Retries = 20
		currentArgs := args
		for attempt := 0; attempt < maxInt32Retries; attempt++ {
			result, err := stmtCtx.ExecContext(ctx, currentArgs)
			if err == nil {
				return result, nil
			}
			if !isInt64ToInt32ConversionError(err) {
				return nil, err
			}
			fieldName := extractFieldNameFromError(err)
			currentArgs = convertInt64ToInt32Args(currentArgs, fieldName, w.query)
		}
		return nil, fmt.Errorf("Int64 to Int32 conversion retry limit (%d) exceeded", maxInt32Retries)
	}

	// Fallback to non-context exec with retry loop
	currentArgs := args
	const maxInt32Retries = 20
	for attempt := 0; attempt < maxInt32Retries; attempt++ {
		values := make([]driver.Value, len(currentArgs))
		for i, arg := range currentArgs {
			values[i] = arg.Value
		}
		result, err := w.stmt.Exec(values)
		if err == nil {
			return result, nil
		}
		if !isInt64ToInt32ConversionError(err) {
			return nil, err
		}
		fieldName := extractFieldNameFromError(err)
		currentArgs = convertInt64ToInt32Args(currentArgs, fieldName, w.query)
	}
	return nil, fmt.Errorf("Int64 to Int32 conversion retry limit (%d) exceeded", maxInt32Retries)
}

// Close closes the underlying statement
func (w *ydbStmtWrapper) Close() error {
	return w.stmt.Close()
}

// NumInput returns the number of placeholder parameters
func (w *ydbStmtWrapper) NumInput() int {
	return w.stmt.NumInput()
}

// Exec implements driver.Stmt interface
func (w *ydbStmtWrapper) Exec(args []driver.Value) (driver.Result, error) {
	currentValues := args
	const maxInt32Retries = 20
	for attempt := 0; attempt < maxInt32Retries; attempt++ {
		result, err := w.stmt.Exec(currentValues)
		if err == nil {
			return result, nil
		}
		if !isInt64ToInt32ConversionError(err) {
			return nil, err
		}
		fieldName := extractFieldNameFromError(err)
		namedArgs := make([]driver.NamedValue, len(currentValues))
		for i, val := range currentValues {
			namedArgs[i] = driver.NamedValue{Ordinal: i + 1, Value: val}
		}
		convertedArgs := convertInt64ToInt32Args(namedArgs, fieldName, w.query)
		currentValues = make([]driver.Value, len(convertedArgs))
		for i, arg := range convertedArgs {
			currentValues[i] = arg.Value
		}
	}
	return nil, fmt.Errorf("Int64 to Int32 conversion retry limit (%d) exceeded", maxInt32Retries)
}

// Query implements driver.Stmt interface
func (w *ydbStmtWrapper) Query(args []driver.Value) (driver.Rows, error) {
	currentValues := args
	const maxInt32Retries = 20
	for attempt := 0; attempt < maxInt32Retries; attempt++ {
		rows, err := w.stmt.Query(currentValues)
		if err == nil {
			return rows, nil
		}
		if !isInt64ToInt32ConversionError(err) {
			return nil, err
		}
		fieldName := extractFieldNameFromError(err)
		namedArgs := make([]driver.NamedValue, len(currentValues))
		for i, val := range currentValues {
			namedArgs[i] = driver.NamedValue{Ordinal: i + 1, Value: val}
		}
		convertedArgs := convertInt64ToInt32Args(namedArgs, fieldName, w.query)
		currentValues = make([]driver.Value, len(convertedArgs))
		for i, arg := range convertedArgs {
			currentValues[i] = arg.Value
		}
	}
	return nil, fmt.Errorf("Int64 to Int32 conversion retry limit (%d) exceeded", maxInt32Retries)
}

func (db *ydbDialect) Init(d *core.DB, uri *core.Uri, drivername, dataSource string) error {
	ydbDriver, err := ydb.Open(context.Background(), dataSource,
		yc.WithInternalCA(),
		// yc.WithCredentials(),
	)
	if err != nil {
		return fmt.Errorf("failed to connect by data source name '%s': %w", dataSource, err)
	}

	connector, err := ydb.Connector(ydbDriver,
		ydb.WithQueryService(true),
		ydb.WithFakeTx(ydb.QueryExecuteQueryMode),
		ydb.WithNumericArgs(),
	)
	if err != nil {
		_ = ydbDriver.Close(context.Background())
		return err
	}

	// Wrap connector to intercept and convert time.Duration arguments
	wrappedConnector := &ydbConnectorWrapper{connector: connector}

	sqldb := sql.OpenDB(wrappedConnector)

	d.DB = sqldb

	return db.Base.Init(core.FromDB(sqldb), db, uri, drivername, dataSource)
}

func (db *ydbDialect) GetIndexes(tableName string) (map[string]*core.Index, error) {
	// YDB index introspection not implemented; return empty so xorm schema operations don't fail.
	return map[string]*core.Index{}, nil
}

func (db *ydbDialect) IndexCheckSql(tableName, idxName string) (string, []interface{}) {
	return "SELECT Path FROM `.sys/partition_stats` where Path LIKE '%/'" +
		" || $1 || '/' || $2 || '/indexImplTable'", []any{tableName, idxName}
}

func (db *ydbDialect) SupportInsertMany() bool {
	return true // TODO:
}

func (db *ydbDialect) SupportCharset() bool {
	return false
}

func (db *ydbDialect) SupportEngine() bool {
	return false
}

func (db *ydbDialect) WithConn(ctx context.Context, f func(context.Context, *sql.Conn) error) error {
	cc, err := db.DB().Conn(ctx)
	if err != nil {
		return err
	}
	defer cc.Close()

	return f(ctx, cc)
}

func (db *ydbDialect) WithConnRaw(ctx context.Context, f func(d interface{}) error) error {
	return db.WithConn(ctx, func(ctx context.Context, cc *sql.Conn) error {
		return cc.Raw(f)
	})
}

func (db *ydbDialect) SetParams(tableParams map[string]string) {
	db.tableParams = tableParams
}

func (db *ydbDialect) IsTableExist(
	ctx context.Context,
	tableName string,
) (_ bool, err error) {
	var exists bool
	err = db.WithConnRaw(ctx, func(dc interface{}) error {
		q, ok := dc.(interface {
			IsTableExists(context.Context, string) (bool, error)
		})
		if !ok {
			return fmt.Errorf("driver hasn't method IsTableExists()")
		}
		exists, err = q.IsTableExists(ctx, tableName)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (db *ydbDialect) TableCheckSql(tableName string) (string, []any) {
	return "SELECT Path FROM `.sys/partition_stats` where Path LIKE '%/' || $1", []any{tableName}
}

func (db *ydbDialect) AutoIncrStr() string {
	return ""
}

func (db *ydbDialect) IsReserved(name string) bool {
	_, ok := ydbReservedWords[strings.ToUpper(name)]
	return ok
}

func (db *ydbDialect) SqlType(column *core.Column) string {
	return toYQLDataType(column.SQLType.Name, column.IsAutoIncrement)
}

// sqlTypeForCreateTable returns the YQL type for a column when creating a table.
// migration_log.timestamp must be DATETIME so INSERT/RETURNING does not fail with "Timestamp to Optional<Date>".
func (db *ydbDialect) sqlTypeForCreateTable(tableName string, col *core.Column) string {
	unquoted := strings.Trim(tableName, "`\"")
	if strings.HasSuffix(unquoted, "migration_log") && col.Name == "timestamp" {
		return yql_DateTime
	}
	return db.SqlType(col)
}

// CreateTableSql implements the Dialect interface; used by Sync2 (e.g. migration_log).
// We override so migration_log.timestamp gets DATETIME (avoids "Timestamp to Optional<Date>" on INSERT).
func (db *ydbDialect) CreateTableSql(table *core.Table, tableName, storeEngine, charset string) string {
	if tableName == "" {
		tableName = table.Name
	}
	parts := make([]string, 0, len(table.ColumnsSeq())+2)
	for _, colName := range table.ColumnsSeq() {
		col := table.GetColumn(colName)
		dataType := db.sqlTypeForCreateTable(tableName, col)
		part := db.Quote(col.Name) + " " + dataType + " "
		if col.IsPrimaryKey && len(table.PrimaryKeys) == 1 {
			part += "PRIMARY KEY " + db.AutoIncrStr() + " "
		}
		if col.Default != "" {
			part += "DEFAULT " + col.Default + " "
		}
		if db.ShowCreateNull() {
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
			quoted[i] = db.Quote(pk)
		}
		parts = append(parts, "PRIMARY KEY ( "+strings.Join(quoted, ", ")+" )")
	}
	return "CREATE TABLE IF NOT EXISTS " + db.Quote(tableName) + " (" + strings.Join(parts, ", ") + ")"
}

// https://pkg.go.dev/database/sql#ColumnType.DatabaseTypeName
func (db *ydbDialect) ColumnTypeKind(t string) int {
	switch t {
	// case "BOOL":
	// 	return core.BOOL_TYPE
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

func (db *ydbDialect) Quote(name string) string {
	return "`" + name + "`" // TODO:
}

func (db *ydbDialect) AddColumnSQL(tableName string, col *core.Column) string {
	tableName = db.Quote(tableName)
	columnName := db.Quote(col.Name)
	dataType := db.SqlType(col)

	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s;", tableName, columnName, dataType))

	return buf.String()
}

// YDB does not support this operation
func (db *ydbDialect) ModifyColumnSQL(tableName string, column *core.Column) string {
	return ""
}

func (db *ydbDialect) DropIndexSql(tableName string, index *core.Index) string {
	tableName = db.Quote(tableName)
	indexName := db.Quote(index.Name)

	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("ALTER TABLE %s DROP INDEX %s;", tableName, indexName))

	return buf.String()
}

func (db *ydbDialect) IndexOnTable() bool {
	return true // TODO:
}

// TODO:
func (db *ydbDialect) IsColumnExist(
	tableName,
	columnName string,
) (_ bool, err error) {
	var exists bool
	ctx := context.TODO()
	err = db.WithConnRaw(ctx, func(dc interface{}) error {
		q, ok := dc.(interface {
			IsColumnExists(context.Context, string, string) (bool, error)
		})
		if !ok {
			return fmt.Errorf("conn hasn't method IsColumnExists()")
		}
		exists, err = q.IsColumnExists(ctx, tableName, columnName)
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return false, err
	}
	return exists, nil
}

// TODO:
func (db *ydbDialect) GetColumns(tableName string) (
	_ []string,
	_ map[string]*core.Column,
	err error,
) {
	ctx := context.TODO()

	colNames := make([]string, 0)
	colMaps := make(map[string]*core.Column)

	// db.nativeDriver

	err = db.WithConnRaw(ctx, func(dc interface{}) error {
		q, ok := dc.(interface {
			GetColumns(context.Context, string) ([]string, error)
			GetColumnType(context.Context, string, string) (string, error)
			IsPrimaryKey(context.Context, string, string) (bool, error)
		})
		if !ok {
			return fmt.Errorf("driver does not support method [GetColumns]")
		}

		colNames, err = q.GetColumns(ctx, tableName)
		if err != nil {
			return err
		}

		for _, colName := range colNames {
			dataType, err := q.GetColumnType(ctx, tableName, colName)
			if err != nil {
				return err
			}
			dataType = removeOptional(dataType)
			isPK, err := q.IsPrimaryKey(ctx, tableName, colName)
			if err != nil {
				return err
			}
			col := &core.Column{
				Name:         colName,
				TableName:    tableName,
				SQLType:      yqlToSQLType(dataType),
				IsPrimaryKey: isPK,
				Nullable:     !isPK,
				Indexes:      make(map[string]int),
			}
			if dataType == "SERIAL" || dataType == "BIGSERIAL" {
				col.IsAutoIncrement = true
			}
			colMaps[colName] = col
		}
		return nil
	})
	if err != nil {
		return nil, nil, err
	}

	return colNames, colMaps, nil
}

func (db *ydbDialect) GetTables() (_ []*core.Table, err error) {
	tables := make([]*core.Table, 0)
	ctx := context.TODO()
	err = db.WithConnRaw(ctx, func(dc interface{}) error {
		q, ok := dc.(interface {
			GetTables(context.Context, string, bool, bool) ([]string, error)
		})
		if !ok {
			return fmt.Errorf("driver does not support method [GetTables]")
		}
		tableNames, err := q.GetTables(ctx, ".", true, true)
		if err != nil {
			return err
		}
		for _, tableName := range tableNames {
			table := core.NewEmptyTable()
			table.Name = tableName
			tables = append(tables, table)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return tables, nil
}

// !datbeohbbh! CreateTableSQL generate `CREATE TABLE` YQL.
// Method does not generate YQL for creating index.
func (db *ydbDialect) CreateTableSQL(
	ctx context.Context,
	_ any,
	table *core.Table,
	tableName string,
) (string, bool, error) {
	tableName = db.Quote(tableName)

	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("CREATE TABLE %s ( ", tableName))

	// 	build primary key
	if len(table.PrimaryKeys) == 0 {
		return "", false, errors.New("table must have at least one primary key")
	}
	pk := make([]string, len(table.PrimaryKeys))
	pkMap := make(map[string]bool)
	for i := 0; i < len(table.PrimaryKeys); i++ {
		pk[i] = db.Quote(table.PrimaryKeys[i])
		pkMap[pk[i]] = true
	}
	primaryKey := fmt.Sprintf("PRIMARY KEY ( %s )", strings.Join(pk, ", "))

	// build column
	columnsList := []string{}
	unquotedTableName := strings.Trim(tableName, "`\"")
	for _, c := range table.Columns() {
		columnName := db.Quote(c.Name)
		dataType := db.SqlType(c)
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

	if len(db.tableParams) > 0 {
		params := make([]string, 0)
		for param, value := range db.tableParams {
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

func (db *ydbDialect) DropTableSQL(tableName string) (string, bool) {
	tableName = db.Quote(tableName)

	var buf strings.Builder
	buf.WriteString(fmt.Sprintf("DROP TABLE %s;", tableName))

	return buf.String(), false
}

type ydbSeqFilter struct {
	Prefix string
	Start  int
}

// TODO:
func (db *ydbDialect) Filters() []core.Filter {
	return []core.Filter{&core.IdFilter{}, &core.SeqFilter{Prefix: "$", Start: 1}}
}

const (
	ydb_grpc_Canceled           uint32 = 1
	ydb_grpc_Unknown            uint32 = 2
	ydb_grpc_InvalidArgument    uint32 = 3
	ydb_grpc_DeadlineExceeded   uint32 = 4
	ydb_grpc_NotFound           uint32 = 5
	ydb_grpc_AlreadyExists      uint32 = 6
	ydb_grpc_PermissionDenied   uint32 = 7
	ydb_grpc_ResourceExhausted  uint32 = 8
	ydb_grpc_FailedPrecondition uint32 = 9
	ydb_grpc_Aborted            uint32 = 10
	ydb_grpc_OutOfRange         uint32 = 11
	ydb_grpc_Unimplemented      uint32 = 12
	ydb_grpc_Internal           uint32 = 13
	ydb_grpc_Unavailable        uint32 = 14
	ydb_grpc_DataLoss           uint32 = 15
	ydb_grpc_Unauthenticated    uint32 = 16
)

const (
	ydb_STATUS_CODE_UNSPECIFIED int32 = 0
	ydb_SUCCESS                 int32 = 400000
	ydb_BAD_REQUEST             int32 = 400010
	ydb_UNAUTHORIZED            int32 = 400020
	ydb_INTERNAL_ERROR          int32 = 400030
	ydb_ABORTED                 int32 = 400040
	ydb_UNAVAILABLE             int32 = 400050
	ydb_OVERLOADED              int32 = 400060
	ydb_SCHEME_ERROR            int32 = 400070
	ydb_GENERIC_ERROR           int32 = 400080
	ydb_TIMEOUT                 int32 = 400090
	ydb_BAD_SESSION             int32 = 400100
	ydb_PRECONDITION_FAILED     int32 = 400120
	ydb_ALREADY_EXISTS          int32 = 400130
	ydb_NOT_FOUND               int32 = 400140
	ydb_SESSION_EXPIRED         int32 = 400150
	ydb_CANCELLED               int32 = 400160
	ydb_UNDETERMINED            int32 = 400170
	ydb_UNSUPPORTED             int32 = 400180
	ydb_SESSION_BUSY            int32 = 400190
)

// https://github.com/ydb-platform/ydb-go-sdk/blob/ca13feb3ca560ac7385e79d4365ffe0cd8c23e21/errors.go#L27
func (db *ydbDialect) IsRetryable(err error) bool {
	var target interface {
		error
		Code() int32
		Name() string
	}
	if errors.Is(err, fmt.Errorf("unknown error")) ||
		errors.Is(err, context.DeadlineExceeded) ||
		errors.Is(err, context.Canceled) {
		return false
	}
	if !errors.As(err, &target) {
		return false
	}

	switch target.Code() {
	case
		int32(ydb_grpc_Unknown),
		int32(ydb_grpc_InvalidArgument),
		int32(ydb_grpc_DeadlineExceeded),
		int32(ydb_grpc_NotFound),
		int32(ydb_grpc_AlreadyExists),
		int32(ydb_grpc_PermissionDenied),
		int32(ydb_grpc_FailedPrecondition),
		int32(ydb_grpc_OutOfRange),
		int32(ydb_grpc_Unimplemented),
		int32(ydb_grpc_DataLoss),
		int32(ydb_grpc_Unauthenticated):
		return false
	case
		int32(ydb_grpc_Canceled),
		int32(ydb_grpc_ResourceExhausted),
		int32(ydb_grpc_Aborted),
		int32(ydb_grpc_Internal),
		int32(ydb_grpc_Unavailable):
		return true
	case
		ydb_STATUS_CODE_UNSPECIFIED,
		ydb_BAD_REQUEST,
		ydb_UNAUTHORIZED,
		ydb_INTERNAL_ERROR,
		ydb_SCHEME_ERROR,
		ydb_GENERIC_ERROR,
		ydb_TIMEOUT,
		ydb_PRECONDITION_FAILED,
		ydb_ALREADY_EXISTS,
		ydb_NOT_FOUND,
		ydb_SESSION_EXPIRED,
		ydb_CANCELLED,
		ydb_UNSUPPORTED:
		return false
	case
		ydb_ABORTED,
		ydb_UNAVAILABLE,
		ydb_OVERLOADED,
		ydb_BAD_SESSION,
		ydb_UNDETERMINED,
		ydb_SESSION_BUSY:
		return true
	default:
		return false
	}
}

type ydbDriver struct {
	core.Base
}

// DSN format: https://github.com/ydb-platform/ydb-go-sdk/blob/a804c31be0d3c44dfd7b21ed49d863619217b11d/connection.go#L339
func (ydbDrv *ydbDriver) Parse(driverName, dataSourceName string) (*core.Uri, error) {
	info := &core.Uri{DbType: core.YDB}

	uri, err := url.Parse(dataSourceName)
	if err != nil {
		return nil, fmt.Errorf("failed on parse data source %v", dataSourceName)
	}

	const (
		secure   = "grpcs"
		insecure = "grpc"
	)

	if uri.Scheme != secure && uri.Scheme != insecure {
		return nil, fmt.Errorf("unsupported scheme %v", uri.Scheme)
	}

	info.Host = uri.Host
	if spl := strings.Split(uri.Host, ":"); len(spl) > 1 {
		info.Host = spl[0]
		info.Port = spl[1]
	}

	info.DbName = uri.Path
	if info.DbName == "" {
		return nil, errors.New("database path can not be empty")
	}

	if uri.User != nil {
		info.Passwd, _ = uri.User.Password()
		info.User = uri.User.Username()
	}

	return info, nil
}

func (db *ydbDialect) RetryOnError(err error) bool {
	return retry.Check(err).MustRetry(true)
}
