package xorm

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"log"
	"net/url"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/ydb-platform/ydb-go-sdk/v3"
	"github.com/ydb-platform/ydb-go-sdk/v3/pkg/xerrors"
	"github.com/ydb-platform/ydb-go-sdk/v3/retry"
	yc "github.com/ydb-platform/ydb-go-yc-metadata"

	"xorm.io/core"
)

const YDB core.DbType = "ydb"

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
// inArgMap maps original positional args to rewritten args (collapsing IN (?,?,...) into one list arg).
type inArgMap struct {
	// segments: in order of placeholders in the rewritten query; each is either one original index or a range to collapse into a list
	segments []inArgSegment
}

type inArgSegment struct {
	singleIdx int // if listCount == 0: use args[singleIdx]
	listStart int // if listCount > 0: use args[listStart:listStart+listCount] as one list arg
	listCount int
}

const inClauseCollapseThreshold = 50 // collapse IN (?,...,?) into IN ? (single list param) when param count >= this

func isSpace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}

func isWordChar(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_'
}

// rewriteQueryInClauses finds large " IN ($k,...,$m)" or " IN (?,...,?)" and replaces with " IN $k" / " IN ?" (single list param).
func rewriteQueryInClauses(query string) (newQuery string, argMap *inArgMap) {
	// Quick rejection: avoid full scan if query cannot contain " IN (...)"
	if !strings.Contains(query, " IN ") && !strings.Contains(query, "IN(") {
		return query, nil
	}
	if strings.Contains(query, "$") {
		if rewritten, m := rewriteQueryInClausesNumeric(query); m != nil {
			return rewritten, m
		}
	}
	return rewriteQueryInClausesQuestion(query)
}

// rewriteQueryInClausesNumeric finds " IN ($k,$k+1,...,$m)" with consecutive ordinals and count >= threshold.
func rewriteQueryInClausesNumeric(query string) (newQuery string, argMap *inArgMap) {
	i := 0
	for i < len(query) {
		j := i
		for j < len(query) && isSpace(query[j]) {
			j++
		}
		// Match "IN" (word boundary) then optional whitespace then "(" (so " IN (", " IN\n(", "IN(")
		if j+2 <= len(query) && (j == 0 || !isWordChar(query[j-1])) && strings.EqualFold(query[j:j+2], "in") {
			k := j + 2
			for k < len(query) && isSpace(query[k]) {
				k++
			}
			if k < len(query) && query[k] == '(' {
				parenStart := k
				parenEnd := findMatchingParen(query, parenStart)
				if parenEnd > parenStart {
					segment := query[parenStart+1 : parenEnd]
					firstNum, count := parseNumericPlaceholderList(segment)
					if count >= inClauseCollapseThreshold && firstNum >= 1 {
						// Replace IN ($k,...,$k+count-1) with IN $k (single list param)
						before := query[:j]
						afterOrig := query[parenEnd+1:]
						// Max ordinal: only scan "after" part, not the long IN list
						lastInList := firstNum + count - 1
						maxAfter := countNumericPlaceholders(afterOrig)
						maxOrd := lastInList
						if maxAfter > maxOrd {
							maxOrd = maxAfter
						}
						after := renumberPlaceholdersAfter(afterOrig, lastInList, -(count - 1))
						newQuery = before + "IN $" + strconv.Itoa(firstNum) + after
						segments := make([]inArgSegment, 0, maxOrd-count+2)
						for a := 0; a < firstNum-1; a++ {
							segments = append(segments, inArgSegment{singleIdx: a})
						}
						segments = append(segments, inArgSegment{listStart: firstNum - 1, listCount: count})
						for a := firstNum + count - 1; a < maxOrd; a++ {
							segments = append(segments, inArgSegment{singleIdx: a})
						}
						return newQuery, &inArgMap{segments: segments}
					}
					// Skip past this IN (...) so we don't re-scan it
					i = parenEnd
				}
			}
		}
		i++
	}
	return query, nil
}

// parseNumericPlaceholderList parses "$n,$n+1,$n+2,..." (with optional newlines/spaces between) and returns (firstNum, count). Returns (0,0) if not consecutive. No slice allocation.
func parseNumericPlaceholderList(segment string) (firstNum int, count int) {
	segment = strings.TrimSpace(segment)
	if segment == "" {
		return 0, 0
	}
	segment = strings.TrimLeft(segment, " \t\n\r,")
	if segment == "" || segment[0] != '$' {
		return 0, 0
	}
	end := 1
	for end < len(segment) && segment[end] >= '0' && segment[end] <= '9' {
		end++
	}
	if end == 1 {
		return 0, 0
	}
	firstNum, _ = strconv.Atoi(segment[1:end])
	count = 1
	segment = segment[end:]
	for len(segment) > 0 {
		segment = strings.TrimLeft(segment, " \t\n\r,")
		if len(segment) == 0 {
			break
		}
		if segment[0] != '$' {
			return 0, 0
		}
		end = 1
		for end < len(segment) && segment[end] >= '0' && segment[end] <= '9' {
			end++
		}
		if end == 1 {
			return 0, 0
		}
		n, _ := strconv.Atoi(segment[1:end])
		if n != firstNum+count {
			return 0, 0
		}
		count++
		segment = segment[end:]
	}
	return firstNum, count
}

// countNumericPlaceholders returns the highest $n ordinal in the query (number of args).
func countNumericPlaceholders(query string) int {
	max := 0
	i := 0
	for i < len(query) {
		if query[i] == '$' && i+1 < len(query) && query[i+1] >= '0' && query[i+1] <= '9' {
			j := i + 1
			for j < len(query) && query[j] >= '0' && query[j] <= '9' {
				j++
			}
			n, _ := strconv.Atoi(query[i+1 : j])
			if n > max {
				max = n
			}
			i = j
			continue
		}
		i++
	}
	return max
}

// renumberPlaceholdersAfter rewrites $m in s to $(m+delta) for every m > threshold.
func renumberPlaceholdersAfter(s string, threshold int, delta int) string {
	var buf strings.Builder
	buf.Grow(len(s)) // output length <= input when delta <= 0
	i := 0
	for i < len(s) {
		if s[i] == '$' && i+1 < len(s) && s[i+1] >= '0' && s[i+1] <= '9' {
			j := i + 1
			for j < len(s) && s[j] >= '0' && s[j] <= '9' {
				j++
			}
			n, _ := strconv.Atoi(s[i+1 : j])
			if n > threshold {
				buf.WriteByte('$')
				buf.WriteString(strconv.Itoa(n + delta))
			} else {
				buf.WriteString(s[i:j])
			}
			i = j
			continue
		}
		buf.WriteByte(s[i])
		i++
	}
	return buf.String()
}

// rewriteQueryInClausesQuestion finds " IN (?, ?, ..., ?)" with >= threshold placeholders.
func rewriteQueryInClausesQuestion(query string) (newQuery string, argMap *inArgMap) {
	var segments []inArgSegment
	var buf strings.Builder
	argIdx := 0
	i := 0
	for i < len(query) {
		if (i == 0 || !isWordChar(query[i-1])) && i+4 <= len(query) {
			j := i
			for j < len(query) && isSpace(query[j]) {
				j++
			}
			if j+4 <= len(query) && strings.EqualFold(query[j:j+4], "in (") {
				parenStart := j + 3
				parenEnd := findMatchingParen(query, parenStart)
				if parenEnd > parenStart {
					inCount := countPlaceholders(query[parenStart+1 : parenEnd])
					if inCount >= inClauseCollapseThreshold {
						buf.WriteString(query[i:j])
						buf.WriteString("IN ?")
						segments = append(segments, inArgSegment{listStart: argIdx, listCount: inCount})
						argIdx += inCount
						i = parenEnd + 1
						continue
					}
				}
			}
		}
		if query[i] == '?' {
			buf.WriteByte('?')
			segments = append(segments, inArgSegment{singleIdx: argIdx})
			argIdx++
			i++
			continue
		}
		buf.WriteByte(query[i])
		i++
	}
	if len(segments) == 0 || segmentsEqualQuery(segments, argIdx) {
		return query, nil
	}
	return buf.String(), &inArgMap{segments: segments}
}

func findMatchingParen(s string, open int) int {
	depth := 1
	for i := open + 1; i < len(s); i++ {
		switch s[i] {
		case '(':
			depth++
		case ')':
			depth--
			if depth == 0 {
				return i
			}
		}
	}
	return -1
}

func countPlaceholders(segment string) int {
	n := 0
	for _, r := range segment {
		if r == '?' {
			n++
		}
	}
	return n
}

// segmentsEqualQuery returns true if the segments just represent 1:1 mapping (no list collapse)
func segmentsEqualQuery(segments []inArgSegment, totalArgs int) bool {
	if len(segments) != totalArgs {
		return false
	}
	for i, s := range segments {
		if s.listCount != 0 || s.singleIdx != i {
			return false
		}
	}
	return true
}

// collapsedCount returns total number of original args that were collapsed into list(s) (for logging).
func (m *inArgMap) collapsedCount() int {
	n := 0
	for _, seg := range m.segments {
		if seg.listCount > 0 {
			n += seg.listCount
		}
	}
	return n
}

func (m *inArgMap) applyNamed(args []driver.NamedValue) ([]driver.NamedValue, error) {
	out := make([]driver.NamedValue, 0, len(m.segments))
	for ord, seg := range m.segments {
		if seg.listCount > 0 {
			list := make([]interface{}, 0, seg.listCount)
			end := seg.listStart + seg.listCount
			if end > len(args) {
				return nil, fmt.Errorf("IN clause arg map: need args[%d:%d], have %d", seg.listStart, end, len(args))
			}
			for j := seg.listStart; j < end; j++ {
				list = append(list, args[j].Value)
			}
			out = append(out, driver.NamedValue{Ordinal: ord + 1, Value: list})
		} else {
			if seg.singleIdx >= len(args) {
				return nil, fmt.Errorf("IN clause arg map: need args[%d], have %d", seg.singleIdx, len(args))
			}
			nm := args[seg.singleIdx]
			nm.Ordinal = ord + 1
			out = append(out, nm)
		}
	}
	return out, nil
}

func (m *inArgMap) applyValues(args []driver.Value) ([]driver.Value, error) {
	out := make([]driver.Value, 0, len(m.segments))
	for _, seg := range m.segments {
		if seg.listCount > 0 {
			list := make([]interface{}, 0, seg.listCount)
			end := seg.listStart + seg.listCount
			if end > len(args) {
				return nil, fmt.Errorf("IN clause arg map: need args[%d:%d], have %d", seg.listStart, end, len(args))
			}
			for j := seg.listStart; j < end; j++ {
				list = append(list, args[j])
			}
			out = append(out, list)
		} else {
			if seg.singleIdx >= len(args) {
				return nil, fmt.Errorf("IN clause arg map: need args[%d], have %d", seg.singleIdx, len(args))
			}
			out = append(out, args[seg.singleIdx])
		}
	}
	return out, nil
}

type ydbStmtWrapper struct {
	stmt   driver.Stmt
	query  string
	inArgs *inArgMap
}

const ydbCostBasedOptimizationPragma = `PRAGMA ydb.CostBasedOptimization = "on";` + "\n"

// shouldUseCostBasedOptimization returns true for "heavy" queries that typically benefit from
// ydb.CostBasedOptimization (dashboard/folder list with permission subqueries). Returns false for
// simple lookups where the pragma often produces worse plans.
func shouldUseCostBasedOptimization(query string) bool {
	q := strings.ToUpper(query)
	hasPermission := strings.Contains(q, "PERMISSION")
	hasDashboardOrFolder := strings.Contains(q, "DASHBOARD") || strings.Contains(q, "FOLDER")
	hasOrderOrLimit := strings.Contains(q, "ORDER BY") || strings.Contains(q, "LIMIT")
	// Heavy pattern: permission + dashboard/folder + ordering/limit (dashboard search, folder list)
	if hasPermission && hasDashboardOrFolder && hasOrderOrLimit {
		return true
	}
	// Multiple nested IN (SELECT ...) — complex permission/folder chain
	inSelectCount := 0
	for i := 0; i < len(q)-10; i++ {
		if strings.HasPrefix(q[i:], "IN (SELECT") {
			inSelectCount++
			if inSelectCount >= 2 {
				return true
			}
		}
	}
	return false
}

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

// ydbILIKEToLowerLikeRe matches "column ILIKE $param" (column may be table.col) for rewrite to LOWER(...) LIKE LOWER(...).
var ydbILIKEToLowerLikeRe = regexp.MustCompile(`([a-zA-Z_][a-zA-Z0-9_]*(?:\.[a-zA-Z_][a-zA-Z0-9_]*)?)\s+ILIKE\s+(\$[a-zA-Z0-9]+)`)

func rewriteILIKEToLowerLike(query string) string {
	if !strings.Contains(query, "ILIKE") {
		return query
	}
	return ydbILIKEToLowerLikeRe.ReplaceAllString(query, "LOWER($1) LIKE LOWER($2)")
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
	}

	rv := reflect.ValueOf(nv.Value)
	if rv.Kind() == reflect.Int {
		nv.Value = rv.Int()
	}

	// Delegate to underlying connector if it implements NamedValueChecker
	if checker, ok := w.conn.(driver.NamedValueChecker); ok {
		return checker.CheckNamedValue(nv)
	}

	return nil
}

// Prepare implements driver.Conn interface
func (w *ydbConnWrapper) Prepare(query string) (driver.Stmt, error) {
	log.Printf("[YDB] Prepare(%q)", query)

	rewritten, inArgs := rewriteQueryInClauses(query)
	if inArgs != nil {
		log.Printf("[YDB] IN clause rewritten: %d params collapsed to single list", inArgs.collapsedCount())
	}
	// Do not rewrite ILIKE to LOWER(...) LIKE: YDB has no LOWER() builtin; use native ILIKE
	if shouldUseCostBasedOptimization(rewritten) {
		rewritten = prependYdbPragma(rewritten)
	}
	stmt, err := w.conn.Prepare(rewritten)
	if err != nil {
		return nil, err
	}
	return &ydbStmtWrapper{stmt: stmt, query: query, inArgs: inArgs}, nil
}

// PrepareContext implements driver.ConnPrepareContext interface
func (w *ydbConnWrapper) PrepareContext(ctx context.Context, query string) (driver.Stmt, error) {
	log.Printf("[YDB] PrepareContext(%q)", query)

	if connCtx, ok := w.conn.(driver.ConnPrepareContext); ok {
		rewritten, inArgs := rewriteQueryInClauses(query)
		if inArgs != nil {
			// Log so operator can confirm IN collapse is applied (xorm [SQL] log shows original query)
			log.Printf("[YDB] IN clause rewritten: %d params collapsed to single list", inArgs.collapsedCount())
		}
		// Do not rewrite ILIKE to LOWER(...) LIKE: YDB has no LOWER() builtin
		if shouldUseCostBasedOptimization(rewritten) {
			rewritten = prependYdbPragma(rewritten)
		}
		stmt, err := connCtx.PrepareContext(ctx, rewritten)
		if err != nil {
			return nil, err
		}
		return &ydbStmtWrapper{stmt: stmt, query: query, inArgs: inArgs}, nil
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

// convertQueryAndArgs converts time.Duration arguments to int64
func convertQueryAndArgs(query string, args []driver.NamedValue) (string, []driver.NamedValue) {
	log.Printf("[YDB] convertQueryAndArgs(%q, %+v)", query, args)

	converted := make([]driver.NamedValue, len(args))
	for i, arg := range args {
		converted[i] = arg
		if duration, ok := arg.Value.(time.Duration); ok {
			converted[i].Value = int64(duration)
		}
	}
	return query, converted
}

// QueryContext intercepts query execution and converts time.Duration to int64
func (w *ydbStmtWrapper) QueryContext(ctx context.Context, args []driver.NamedValue) (driver.Rows, error) {
	_, args = convertQueryAndArgs(w.query, args)

	if w.inArgs != nil {
		var err error
		args, err = w.inArgs.applyNamed(args)
		if err != nil {
			return nil, err
		}
	}

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

	var resultRows driver.Rows
	err := retry.Retry(ctx, func(ctx context.Context) (err error) {
		if stmtCtx, ok := w.stmt.(driver.StmtQueryContext); ok {
			const maxInt32Retries = 20
			currentArgs := args
			for attempt := 0; attempt < maxInt32Retries; attempt++ {
				rows, err := stmtCtx.QueryContext(ctx, currentArgs)
				if err == nil {
					resultRows = &rowsWrapper{Rows: rows}
					return nil
				}
				if !isInt64ToInt32ConversionError(err) {
					return err
				}
				fieldName := extractFieldNameFromError(err)
				currentArgs = convertInt64ToInt32Args(currentArgs, fieldName, w.query)
			}
			return fmt.Errorf("Int64 to Int32 conversion retry limit (%d) exceeded", maxInt32Retries)
		}
		currentArgs := args
		const maxInt32Retries = 20
		for attempt := 0; attempt < maxInt32Retries; attempt++ {
			values := make([]driver.Value, len(currentArgs))
			for i, arg := range currentArgs {
				values[i] = arg.Value
			}
			rows, err := w.stmt.Query(values)
			if err == nil {
				resultRows = rows
				return nil
			}
			if !isInt64ToInt32ConversionError(err) {
				return err
			}
			fieldName := extractFieldNameFromError(err)
			currentArgs = convertInt64ToInt32Args(currentArgs, fieldName, w.query)
		}
		return fmt.Errorf("Int64 to Int32 conversion retry limit (%d) exceeded", maxInt32Retries)
	})
	if err != nil {
		return nil, err
	}
	return resultRows, nil
}

// ExecContext intercepts exec execution and converts time.Duration to int64
func (w *ydbStmtWrapper) ExecContext(ctx context.Context, args []driver.NamedValue) (driver.Result, error) {
	_, args = convertQueryAndArgs(w.query, args)

	if w.inArgs != nil {
		var err error
		args, err = w.inArgs.applyNamed(args)
		if err != nil {
			return nil, err
		}
	}

	var result driver.Result
	err := retry.Retry(ctx, func(ctx context.Context) (err error) {
		if stmtCtx, ok := w.stmt.(driver.StmtExecContext); ok {
			const maxInt32Retries = 20
			currentArgs := args
			for attempt := 0; attempt < maxInt32Retries; attempt++ {
				res, err := stmtCtx.ExecContext(ctx, currentArgs)
				if err == nil {
					result = res
					return nil
				}
				if !isInt64ToInt32ConversionError(err) {
					return err
				}
				fieldName := extractFieldNameFromError(err)
				currentArgs = convertInt64ToInt32Args(currentArgs, fieldName, w.query)
			}
			return fmt.Errorf("Int64 to Int32 conversion retry limit (%d) exceeded", maxInt32Retries)
		}
		currentArgs := args
		const maxInt32Retries = 20
		for attempt := 0; attempt < maxInt32Retries; attempt++ {
			values := make([]driver.Value, len(currentArgs))
			for i, arg := range currentArgs {
				values[i] = arg.Value
			}
			res, err := w.stmt.Exec(values)
			if err == nil {
				result = res
				return nil
			}
			if !isInt64ToInt32ConversionError(err) {
				return err
			}
			fieldName := extractFieldNameFromError(err)
			currentArgs = convertInt64ToInt32Args(currentArgs, fieldName, w.query)
		}
		return fmt.Errorf("Int64 to Int32 conversion retry limit (%d) exceeded", maxInt32Retries)
	})
	if err != nil {
		return nil, err
	}
	return result, nil
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
	return nil, xerrors.WithStackTrace(errors.New("[YDB] Exec not supported, use ExecContext instead"))
}

// Query implements driver.Stmt interface
func (w *ydbStmtWrapper) Query(args []driver.Value) (driver.Rows, error) {
	return nil, xerrors.WithStackTrace(errors.New("[YDB] Query not supported, use QueryContext instead"))
}

func (db *ydbDialect) Init(d *core.DB, uri *core.Uri, drivername, dataSource string) error {
	ydbDriver, err := ydb.Unwrap(d.DB)
	if err != nil {
		return fmt.Errorf("unwrapping ydb driver: %w", err)
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
	return "SELECT Path FROM `.sys/partition_stats` WHERE Path LIKE '%/' || $1", []any{tableName}
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

// TODO:
func (db *ydbDialect) Filters() []core.Filter {
	return []core.Filter{&core.IdFilter{}, &core.SeqFilter{Prefix: "$", Start: 1}}
}

func (db *ydbDialect) IsRetryable(err error) bool {
	return retry.Check(err).MustRetry(false)
}

type ydbDriver struct {
	core.Base
}

// DSN format: https://github.com/ydb-platform/ydb-go-sdk/blob/a804c31be0d3c44dfd7b21ed49d863619217b11d/connection.go#L339
func (ydbDrv *ydbDriver) Parse(driverName, dataSourceName string) (*core.Uri, error) {
	info := &core.Uri{DbType: YDB}

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
