package filters

import (
	"github.com/grafana/grafana/pkg/util/xorm/core"
)

var (
	_ core.Filter         = (*Numbers)(nil)
	_ core.FilterWithArgs = (*Numbers)(nil)
)

type Numbers struct{}

func (f *Numbers) Do(sql string, _ core.Dialect, _ *core.Table) string {
	return sql
}

func (f *Numbers) DoWithArgs(sql string, _ core.Dialect, _ *core.Table, args ...any) (string, []any) {
	for i := range args {
		switch v := args[i].(type) {
		case int:
			args[i] = int64(v)
		case int32:
			args[i] = int64(v)
		case float32:
			args[i] = float64(v)
		default:
			// nop
		}
	}

	return sql, args
}
