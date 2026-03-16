package filters

import (
	"time"

	"github.com/grafana/grafana/pkg/util/xorm/core"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

var (
	_ core.Filter         = (*ConvertTimeToDatetime64)(nil)
	_ core.FilterWithArgs = (*ConvertTimeToDatetime64)(nil)
)

type ConvertTimeToDatetime64 struct{}

func (f *ConvertTimeToDatetime64) Do(sql string, _ core.Dialect, _ *core.Table) string {
	panic("unexpected call Do, expected DoWithArgs")
}

func (f *ConvertTimeToDatetime64) DoWithArgs(sql string, _ core.Dialect, _ *core.Table, args ...any) (string, []any) {
	for i := range args {
		switch t := args[i].(type) {
		case time.Time:
			args[i] = types.Datetime64Value(t.Unix())
		case *time.Time:
			args[i] = types.NullableDatetime64ValueFromTime(t)
		default:
			// nop
		}
	}

	return sql, args
}
