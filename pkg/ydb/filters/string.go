package filters

import (
	"time"

	"github.com/grafana/grafana/pkg/util/xorm/core"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

var (
	_ core.Filter         = (*ConvertStringToDatetime64)(nil)
	_ core.FilterWithArgs = (*ConvertStringToDatetime64)(nil)
)

type ConvertStringToDatetime64 struct{}

func (f *ConvertStringToDatetime64) Do(sql string, _ core.Dialect, _ *core.Table) string {
	panic("unexpected call Do, expected DoWithArgs")
}

func (f *ConvertStringToDatetime64) DoWithArgs(sql string, _ core.Dialect, _ *core.Table, args ...any) (string, []any) {
	for i := range args {
		switch v := args[i].(type) {
		case string:
			if t, err := time.Parse(time.DateTime, v); err == nil {
				args[i] = types.Datetime64Value(t.Unix())
			}
		default:
			// nop
		}
	}

	return sql, args
}
