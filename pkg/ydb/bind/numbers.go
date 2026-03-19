package bind

var _ Binder = (*ConvertNumbersToInt64)(nil)

type ConvertNumbersToInt64 struct{}

func (f *ConvertNumbersToInt64) Rebind(sql string, args ...any) (string, []any, error) {
	out := make([]any, len(args))
	for i := range args {
		switch v := args[i].(type) {
		case int:
			out[i] = int64(v)
		case int32:
			out[i] = int64(v)
		case float32:
			out[i] = float64(v)
		default:
			out[i] = args[i]
		}
	}
	return sql, out, nil
}
