package filters

import (
	"time"

	"github.com/grafana/grafana/pkg/util/xorm/core"
)

var (
	_ core.Filter         = (*Duration)(nil)
	_ core.FilterWithArgs = (*Duration)(nil)
)

type Duration struct{}

func (f *Duration) Do(sql string, _ core.Dialect, _ *core.Table) string {
	return sql
}

func (f *Duration) DoWithArgs(sql string, _ core.Dialect, _ *core.Table, args ...any) (string, []any) {
	for i := range args {
		if d, has := args[i].(time.Duration); has {
			// args[i] = int64(d / time.Millisecond)
			args[i] = int64(d)
		}
	}

	return sql, args
}
