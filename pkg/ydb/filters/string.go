package filters

import (
	"time"

	"github.com/grafana/grafana/pkg/util/xorm/core"
	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

var (
	_ core.Filter         = (*String)(nil)
	_ core.FilterWithArgs = (*String)(nil)
)

type String struct{}

func (f *String) Do(sql string, _ core.Dialect, _ *core.Table) string {
	return sql
}

func (f *String) DoWithArgs(sql string, _ core.Dialect, _ *core.Table, args ...any) (string, []any) {
	for i := range args {
		switch v := args[i].(type) {
		case string:
			if t, err := time.Parse(time.DateTime, v); err == nil {
				args[i] = types.Datetime64Value(t.Unix())
				//} else {
				//	if b := xstring.ToBytes(v); len(b) > 0 {
				//		args[i] = b
				//	} else {
				//		args[i] = []byte{}
				//	}
			}
		default:
			// nop
		}
	}

	return sql, args
}
