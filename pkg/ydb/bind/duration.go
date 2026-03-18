package bind

import (
	"database/sql/driver"
	"time"
)

var _ Binder = (*ConvertDurationToInt64)(nil)

type ConvertDurationToInt64 struct{}

func (f *ConvertDurationToInt64) Rebind(sql string, args ...driver.NamedValue) (string, []driver.NamedValue, error) {
	for i := range args {
		if d, has := args[i].Value.(time.Duration); has {
			args[i].Value = int64(d)
		}
	}
	return sql, args, nil
}
