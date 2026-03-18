package bind

import "database/sql/driver"

var _ Binder = (*ConvertNumbersToInt64)(nil)

type ConvertNumbersToInt64 struct{}

func (f *ConvertNumbersToInt64) Rebind(sql string, args ...driver.NamedValue) (string, []driver.NamedValue, error) {
	for i := range args {
		switch v := args[i].Value.(type) {
		case int:
			args[i].Value = int64(v)
		case int32:
			args[i].Value = int64(v)
		case float32:
			args[i].Value = float64(v)
		}
	}
	return sql, args, nil
}
