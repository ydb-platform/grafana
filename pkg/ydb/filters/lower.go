package filters

import (
	"strings"

	"github.com/grafana/grafana/pkg/util/xorm/core"
)

var (
	_ core.Filter = (*Lower)(nil)
)

type Lower struct{}

const lowerFrom = "LOWER"
const lowerTo = "Unicode::ToLower"

func (Lower) Do(sql string, dialect core.Dialect, table *core.Table) string {
	if !strings.Contains(sql, lowerFrom) {
		return sql
	}
	var out strings.Builder
	i := 0
	for i < len(sql) {
		replaced := false
		if i+len(lowerFrom) <= len(sql) && sql[i:i+len(lowerFrom)] == lowerFrom {
			prevOK := i == 0 || !isWordChar(sql[i-1])
			nextOK := i+len(lowerFrom) >= len(sql) || !isWordChar(sql[i+len(lowerFrom)])
			if prevOK && nextOK && i+len(lowerFrom) < len(sql) && sql[i+len(lowerFrom)] == '(' {
				// LOWER(arg) → Unicode::ToLower(CAST(arg AS Text))
				start := i + len(lowerFrom) + 1
				depth := 1
				end := start
			findClose:
				for end < len(sql) {
					switch sql[end] {
					case '(':
						depth++
					case ')':
						depth--
						if depth == 0 {
							out.WriteString("Unicode::ToLower(CAST(")
							out.WriteString(sql[start:end])
							out.WriteString(" AS Text))")
							i = end + 1
							replaced = true
							break findClose
						}
					}
					end++
				}
			}
			if !replaced && prevOK && nextOK {
				out.WriteString(lowerTo)
				i += len(lowerFrom)
				continue
			}
		}
		if replaced {
			continue
		}
		out.WriteByte(sql[i])
		i++
	}
	return out.String()
}
