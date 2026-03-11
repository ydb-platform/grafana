package filters

import (
	"github.com/grafana/grafana/pkg/util/xorm/core"
)

var (
	_ core.Filter         = (*Bool)(nil)
	_ core.FilterWithArgs = (*Bool)(nil)
)

type Bool struct{}

func (f *Bool) Do(sql string, _ core.Dialect, _ *core.Table) string {
	return sql
}

func (f *Bool) DoWithArgs(sql string, _ core.Dialect, _ *core.Table, args ...any) (string, []any) {
	for i := range args {
		switch v := args[i].(type) {
		case bool:
			if v {
				args[i] = uint8(1)
			} else {
				args[i] = uint8(0)
			}
		default:
			// nop
		}
	}

	return sql, args
}
