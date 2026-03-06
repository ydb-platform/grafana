package filters

import (
	"github.com/grafana/grafana/pkg/util/xorm/core"
)

var (
	_ core.Filter         = (*Literal)(nil)
	_ core.FilterWithArgs = (*Literal)(nil)
)

type Literal struct{}

func (f *Literal) Do(sql string, _ core.Dialect, _ *core.Table) string {
	return sql
}

// literalSpan describes one string literal ('...' or "...") to replace with ?.
type literalSpan struct {
	start, end int
	value      string
}

// findCreateTableBodies returns [start,end) ranges of CREATE TABLE (...); bodies.
// Literals inside these ranges are not converted to parameters.
func findCreateTableBodies(sql string) [][2]int {
	var ranges [][2]int
	i := 0
	for i < len(sql) {
		if i+12 <= len(sql) && (sql[i] == 'C' || sql[i] == 'c') &&
			(sql[i+1] == 'R' || sql[i+1] == 'r') && (sql[i+2] == 'E' || sql[i+2] == 'e') &&
			(sql[i+3] == 'A' || sql[i+3] == 'a') && (sql[i+4] == 'T' || sql[i+4] == 't') &&
			(sql[i+5] == 'E' || sql[i+5] == 'e') && (sql[i+6] == ' ' || sql[i+6] == '\t') &&
			(sql[i+7] == 'T' || sql[i+7] == 't') && (sql[i+8] == 'A' || sql[i+8] == 'a') &&
			(sql[i+9] == 'B' || sql[i+9] == 'b') && (sql[i+10] == 'L' || sql[i+10] == 'l') &&
			(sql[i+11] == 'E' || sql[i+11] == 'e') &&
			(i == 0 || !isWordByte(sql[i-1])) {
			afterCreate := i + 12
			// skip IF NOT EXISTS and table name until we see '('
			for afterCreate < len(sql) && sql[afterCreate] != '(' {
				if sql[afterCreate] == '`' {
					afterCreate++
					for afterCreate < len(sql) && sql[afterCreate] != '`' {
						afterCreate++
					}
					if afterCreate < len(sql) {
						afterCreate++
					}
				} else {
					afterCreate++
				}
			}
			if afterCreate < len(sql) && sql[afterCreate] == '(' {
				bodyEnd := findMatchingParen(sql, afterCreate)
				if bodyEnd > afterCreate {
					ranges = append(ranges, [2]int{afterCreate, bodyEnd + 1})
					i = bodyEnd + 1
					continue
				}
			}
		}
		i++
	}
	return ranges
}

func findMatchingParen(sql string, open int) int {
	if open >= len(sql) || sql[open] != '(' {
		return -1
	}
	depth := 1
	inString := false
	var quote byte
	i := open + 1
	for i < len(sql) {
		c := sql[i]
		if inString {
			if c == quote {
				inString = false
			} else if c == '\\' && i+1 < len(sql) {
				i++
			}
			i++
			continue
		}
		if c == '\'' || c == '"' || c == '`' {
			inString = true
			quote = c
			i++
			continue
		}
		if c == '(' {
			depth++
			i++
			continue
		}
		if c == ')' {
			depth--
			if depth == 0 {
				return i
			}
			i++
			continue
		}
		i++
	}
	return -1
}

func inExcludeRange(pos int, exclude [][2]int) bool {
	for _, r := range exclude {
		if pos >= r[0] && pos < r[1] {
			return true
		}
	}
	return false
}

// findLiterals finds all string literals in single or double quotes.
// Backticks are not treated as value literals (used for identifiers). Numeric literals
// are not extracted (LIMIT/OFFSET etc. are handled by other filters).
// Literals inside CREATE TABLE (...); bodies are skipped.
func (f *Literal) findLiterals(sql string) []literalSpan {
	exclude := findCreateTableBodies(sql)
	var out []literalSpan
	i := 0
	for i < len(sql) {
		if end := skipExclude(i, exclude); end > i {
			i = end
			continue
		}
		if c := sql[i]; c == '\'' || c == '"' {
			start := i
			val, end, ok := f.parseQuotedString(sql, i)
			if ok {
				out = append(out, literalSpan{start, end, val})
				i = end
				continue
			}
		}
		i++
	}
	return out
}

func skipExclude(pos int, exclude [][2]int) int {
	for _, r := range exclude {
		if pos >= r[0] && pos < r[1] {
			return r[1]
		}
	}
	return 0
}

func (f *Literal) parseQuotedString(sql string, i int) (string, int, bool) {
	if i >= len(sql) {
		return "", i, false
	}
	quote := sql[i]
	i++
	start := i
	for i < len(sql) {
		c := sql[i]
		if c == quote {
			i++
			return sql[start : i-1], i, true
		}
		if c == '\\' && i+1 < len(sql) {
			i += 2
			continue
		}
		i++
	}
	return "", i, false
}

func (f *Literal) DoWithArgs(sql string, _ core.Dialect, _ *core.Table, args ...any) (string, []any) {
	literals := f.findLiterals(sql)
	if len(literals) == 0 {
		return sql, args
	}

	// Collect all '?' positions (outside string literals).
	var qPositions []int
	inString := false
	var quote byte
	for i := 0; i < len(sql); i++ {
		c := sql[i]
		if inString {
			if c == quote {
				inString = false
			} else if c == '\\' && i+1 < len(sql) {
				i++
			}
			continue
		}
		if c == '\'' || c == '"' {
			inString = true
			quote = c
			continue
		}
		if c == '?' {
			qPositions = append(qPositions, i)
		}
	}

	// Merge ? and literal positions; build new args in SQL order.
	type segment struct {
		pos int
		q   bool
		lit int
	}
	var segs []segment
	for _, pos := range qPositions {
		segs = append(segs, segment{pos: pos, q: true})
	}
	for j := range literals {
		segs = append(segs, segment{pos: literals[j].start, q: false, lit: j})
	}
	for i := 0; i < len(segs); i++ {
		for j := i + 1; j < len(segs); j++ {
			if segs[j].pos < segs[i].pos {
				segs[i], segs[j] = segs[j], segs[i]
			}
		}
	}
	newArgs := make([]any, 0, len(args)+len(literals))
	argIdx := 0
	for _, s := range segs {
		if s.q {
			if argIdx < len(args) {
				newArgs = append(newArgs, args[argIdx])
				argIdx++
			}
		} else {
			newArgs = append(newArgs, literals[s.lit].value)
		}
	}

	// Replace from end to start so positions don't shift.
	for j := len(literals) - 1; j >= 0; j-- {
		lit := literals[j]
		sql = sql[:lit.start] + "?" + sql[lit.end:]
	}

	return sql, newArgs
}
