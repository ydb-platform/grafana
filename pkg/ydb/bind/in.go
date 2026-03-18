package bind

import "database/sql/driver"

var _ Binder = (*ConvertInArgsToList)(nil)

type ConvertInArgsToList struct{}

type ydbInRange struct {
	start, end int
	count      int
	firstQ     int
}

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
		if (c == 'I' || c == 'i') && i+2 <= len(sql) && (sql[i+1] == 'N' || sql[i+1] == 'n') && (i == 0 || !isWordByte(sql[i-1])) {
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

func (f *ConvertInArgsToList) Rebind(sql string, args ...driver.NamedValue) (string, []driver.NamedValue, error) {
	ranges := f.findInClauses(sql)
	if len(ranges) == 0 {
		return sql, args, nil
	}

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

	newArgs := make([]driver.NamedValue, 0, len(args))
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
			slice := make([]any, in.count)
			for k := 0; k < in.count && argIdx+k < len(args); k++ {
				slice[k] = args[argIdx+k].Value
			}
			newArgs = append(newArgs, driver.NamedValue{Value: slice})
			argIdx += in.count
		} else if in == nil {
			if argIdx < len(args) {
				newArgs = append(newArgs, args[argIdx])
			}
			argIdx++
		}
	}

	outSQL := sql
	for i := len(ranges) - 1; i >= 0; i-- {
		r := ranges[i]
		outSQL = outSQL[:r.start] + "IN ?" + outSQL[r.end:]
	}

	return outSQL, newArgs, nil
}
