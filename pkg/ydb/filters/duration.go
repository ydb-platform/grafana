package filters

import (
	"time"

	"github.com/grafana/grafana/pkg/util/xorm/core"
)

var (
	_ core.Filter         = (*ConvertDurationToInt64)(nil)
	_ core.FilterWithArgs = (*ConvertDurationToInt64)(nil)
)

type ConvertDurationToInt64 struct{}

func (f *ConvertDurationToInt64) Do(sql string, _ core.Dialect, _ *core.Table) string {
	panic("unexpected call Do, expected DoWithArgs")
}

func (f *ConvertDurationToInt64) DoWithArgs(sql string, _ core.Dialect, _ *core.Table, args ...any) (string, []any) {
	for i := range args {
		if d, has := args[i].(time.Duration); has {
			// args[i] = int64(d / time.Millisecond)
			args[i] = int64(d)
		}
	}

	return sql, args
}
