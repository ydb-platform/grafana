package filters

import (
	"strings"

	"github.com/grafana/grafana/pkg/util/xorm/core"
)

var (
	_ core.Filter = (*Substr)(nil)
)

type Substr struct{}

const substrFrom = "SUBSTR"
const substrTo = "SUBSTRING"

func (s Substr) Do(sql string, dialect core.Dialect, table *core.Table) string {
	if !strings.Contains(sql, "SUBSTR") {
		return sql
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
	return out.String()
}
