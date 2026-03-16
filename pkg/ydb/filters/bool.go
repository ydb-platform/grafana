package filters

import (
	"github.com/grafana/grafana/pkg/util/xorm/core"
)

var (
	_ core.Filter         = (*ConvertBoolToUint8)(nil)
	_ core.FilterWithArgs = (*ConvertBoolToUint8)(nil)
)

type ConvertBoolToUint8 struct{}

func (f *ConvertBoolToUint8) Do(sql string, _ core.Dialect, _ *core.Table) string {
	panic("unexpected call Do, expected DoWithArgs")
}

func (f *ConvertBoolToUint8) DoWithArgs(sql string, _ core.Dialect, _ *core.Table, args ...any) (string, []any) {
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
