package bind

import (
	"time"

	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

var _ Binder = (*ConvertStringToDatetime64)(nil)

type ConvertStringToDatetime64 struct{}

func (f *ConvertStringToDatetime64) Rebind(sql string, args ...any) (string, []any, error) {
	out := make([]any, len(args))
	for i := range args {
		if v, ok := args[i].(string); ok {
			if t, err := time.Parse(time.DateTime, v); err == nil {
				out[i] = types.Datetime64Value(t.Unix())
			} else {
				out[i] = args[i]
			}
		} else {
			out[i] = args[i]
		}
	}
	return sql, out, nil
}
