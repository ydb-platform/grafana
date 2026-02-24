package filters

import (
	"strings"

	"github.com/grafana/grafana/pkg/util/xorm/core"
)

var (
	_ core.Filter = (*Replace)(nil)
)

type Replace struct{}

func isWordChar(b byte) bool {
	return (b >= 'a' && b <= 'z') || (b >= 'A' && b <= 'Z') || (b >= '0' && b <= '9') || b == '_'
}

const replaceFrom = "REPLACE"
const replaceTo = "Unicode::ReplaceAll"

func (Replace) Do(sql string, dialect core.Dialect, table *core.Table) string {
	if !strings.Contains(sql, "REPLACE") {
		return sql
	}
	var out strings.Builder
	i := 0
	for i < len(sql) {
		if i+len(replaceFrom) <= len(sql) && sql[i:i+len(replaceFrom)] == replaceFrom {
			prevOK := i == 0 || !isWordChar(sql[i-1])
			nextOK := i+len(replaceFrom) >= len(sql) || !isWordChar(sql[i+len(replaceFrom)])
			if prevOK && nextOK {
				out.WriteString(replaceTo)
				i += len(replaceFrom)
				continue
			}
		}
		out.WriteByte(sql[i])
		i++
	}
	return out.String()
}
