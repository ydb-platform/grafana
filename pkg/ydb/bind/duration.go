package bind

import "time"

var _ Binder = (*ConvertDurationToInt64)(nil)

type ConvertDurationToInt64 struct{}

func (f *ConvertDurationToInt64) Rebind(sql string, args ...any) (string, []any, error) {
	out := make([]any, len(args))
	for i := range args {
		if d, ok := args[i].(time.Duration); ok {
			out[i] = int64(d)
		} else {
			out[i] = args[i]
		}
	}
	return sql, out, nil
}
