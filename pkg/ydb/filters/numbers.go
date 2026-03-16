package filters

import (
	"github.com/grafana/grafana/pkg/util/xorm/core"
)

var (
	_ core.Filter         = (*ConvertNumbersToInt64)(nil)
	_ core.FilterWithArgs = (*ConvertNumbersToInt64)(nil)
)

type ConvertNumbersToInt64 struct{}

func (f *ConvertNumbersToInt64) Do(sql string, _ core.Dialect, _ *core.Table) string {
	panic("unexpected call Do, expected DoWithArgs")
}

func (f *ConvertNumbersToInt64) DoWithArgs(sql string, _ core.Dialect, _ *core.Table, args ...any) (string, []any) {
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
