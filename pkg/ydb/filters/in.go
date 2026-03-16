package filters

import "github.com/grafana/grafana/pkg/util/xorm/core"

var (
	_ core.Filter         = (*ConvertInArgsToList)(nil)
	_ core.FilterWithArgs = (*ConvertInArgsToList)(nil)
)

// IN filter SQL replace IN (?, ?, ?) to IN ? and many args (1,2,3) replace to single arg []any{1,2,3}
type ConvertInArgsToList struct{}

// ydbInRange describes one "IN (?,...,?)" clause in the SQL (only placeholders, no subquery).
type ydbInRange struct {
	start, end int
	count      int
	firstQ     int
}

// findInClauses finds all "IN (?,?,?)" style clauses (no regex): IN, optional space, (, only ? , space, ).
func (f *ConvertInArgsToList) findInClauses(sql string) []ydbInRange {
	var ranges []ydbInRange
	inString := false
	var quote byte
	i := 0
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
		if c == '\'' || c == '"' {
			inString = true
			quote = c
			i++
			continue
		}
		// Look for IN not preceded by word char
		if (c == 'I' || c == 'i') && i+2 <= len(sql) && (sql[i+1] == 'N' || sql[i+1] == 'n') && (i == 0 || !f.isWordByte(sql[i-1])) {
			start := i
			i += 2
			for i < len(sql) && (sql[i] == ' ' || sql[i] == '\t' || sql[i] == '\n' || sql[i] == '\r') {
				i++
			}
			if i < len(sql) && sql[i] == '(' {
				i++
				count := 0
				firstQ := -1
				ok := true
				for i < len(sql) {
					b := sql[i]
					if b == '?' {
						if firstQ < 0 {
							firstQ = i
						}
						count++
						i++
						continue
					}
					if b == ' ' || b == '\t' || b == ',' || b == '\n' || b == '\r' {
						i++
						continue
					}
					if b == ')' {
						i++
						if count > 0 && ok {
							ranges = append(ranges, ydbInRange{start, i, count, firstQ})
						}
						break
					}
					ok = false
					i++
				}
				continue
			}
		}
		i++
	}
	return ranges
}

func (f *ConvertInArgsToList) isWordByte(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_'
}

func (f *ConvertInArgsToList) Do(sql string, _ core.Dialect, _ *core.Table) string {
	panic("unexpected call Do, expected DoWithArgs")
}

func (f *ConvertInArgsToList) DoWithArgs(sql string, _ core.Dialect, _ *core.Table, args ...any) (string, []any) {
	ranges := f.findInClauses(sql)
	if len(ranges) == 0 {
		return sql, args
	}

	// Collect all '?' positions in SQL in order (skip ? inside string literals).
	var qPositions []int
	inString := false
	var quote byte
	for i := 0; i < len(sql); i++ {
		c := sql[i]
		if inString {
			if c == quote {
				inString = false
			} else if c == '\\' && i+1 < len(sql) {
				i++ // skip escaped char
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

	// Build new args: for each '?' in order, either append one arg or append one slice and skip next (count-1) args.
	newArgs := make([]any, 0, len(args))
	argIdx := 0
	for _, pos := range qPositions {
		var in *ydbInRange
		var isFirst bool
		for j := range ranges {
			r := &ranges[j]
			if pos >= r.start && pos < r.end {
				in = r
				isFirst = (pos == r.firstQ)
				break
			}
		}
		if in != nil && isFirst {
			newArgs = append(newArgs, append([]any(nil), args[argIdx:argIdx+in.count]...))
			argIdx += in.count
		} else if in == nil {
			newArgs = append(newArgs, args[argIdx])
			argIdx++
		}
	}

	// Replace every IN (?,...,?) with "IN ?" in SQL (from end to start so indices stay valid).
	outSQL := sql
	for i := len(ranges) - 1; i >= 0; i-- {
		r := ranges[i]
		outSQL = outSQL[:r.start] + "IN ?" + outSQL[r.end:]
	}

	return outSQL, newArgs
}
