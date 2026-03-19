package bind

import (
	"time"

	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

var _ Binder = (*ConvertTimeToDatetime64)(nil)

type ConvertTimeToDatetime64 struct{}

func (f *ConvertTimeToDatetime64) Rebind(sql string, args ...any) (string, []any, error) {
	out := make([]any, len(args))
	for i := range args {
		switch t := args[i].(type) {
		case time.Time:
			out[i] = types.Datetime64Value(t.Unix())
		case *time.Time:
			out[i] = types.NullableDatetime64ValueFromTime(t)
		default:
			out[i] = args[i]
		}
	}
	return sql, out, nil
}
