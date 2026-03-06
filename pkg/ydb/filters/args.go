package filters

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/grafana/grafana/pkg/util/xorm/core"
)

var (
	_ core.Filter         = (*Args)(nil)
	_ core.FilterWithArgs = (*Args)(nil)
)

// Args filter SQL replace ?, ? ... to $p1, $p2 ...
type Args struct{}

func (f *Args) Do(sql string, _ core.Dialect, _ *core.Table) string {
	return sql
}

func (f *Args) DoWithArgs(query string, _ core.Dialect, _ *core.Table, args ...any) (string, []any) {
	var buf strings.Builder
	var beginQuote rune
	index := 0
	for _, c := range query {
		if beginQuote == 0 {
			switch c {
			case '\'', '"':
				beginQuote = c
				buf.WriteRune(c)
			case '?':
				paramName := fmt.Sprintf("p%d", index+1)
				buf.WriteString("$" + paramName)
				args[index] = sql.Named(paramName, args[index])
				index++
			default:
				buf.WriteRune(c)
			}
		} else {
			switch c {
			case beginQuote:
				beginQuote = 0
				buf.WriteRune(c)
			default:
				buf.WriteRune(c)
			}
		}
	}
	return buf.String(), args
}
