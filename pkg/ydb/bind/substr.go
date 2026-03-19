package bind

import "strings"

var _ Binder = (*ConvertSubstrToSubstring)(nil)

type ConvertSubstrToSubstring struct{}

const substrFrom = "SUBSTR"
const substrTo = "SUBSTRING"

func (ConvertSubstrToSubstring) Rebind(sql string, args ...any) (string, []any, error) {
	if !strings.Contains(sql, "SUBSTR") {
		return sql, args, nil
	}
	var out strings.Builder
	i := 0
	for i < len(sql) {
		if sql[i] == '"' || sql[i] == '\'' {
			quote := sql[i]
			out.WriteByte(sql[i])
			i++
			for i < len(sql) {
				c := sql[i]
				if c == quote {
					out.WriteByte(sql[i])
					i++
					break
				}
				if c == '\\' && i+1 < len(sql) {
					out.WriteByte(sql[i])
					out.WriteByte(sql[i+1])
					i += 2
					continue
				}
				out.WriteByte(sql[i])
				i++
			}
			continue
		}
		if i+len(substrFrom) <= len(sql) && sql[i:i+len(substrFrom)] == substrFrom {
			prevOK := i == 0 || !isWordChar(sql[i-1])
			nextOK := i+len(substrFrom) >= len(sql) || !isWordChar(sql[i+len(substrFrom)])
			if prevOK && nextOK {
				out.WriteString(substrTo)
				i += len(substrFrom)
				continue
			}
		}
		out.WriteByte(sql[i])
		i++
	}
	return out.String(), args, nil
}
