package bind

import (
	"database/sql/driver"
	"time"

	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

var _ Binder = (*ConvertTimeToDatetime64)(nil)

type ConvertTimeToDatetime64 struct{}

func (f *ConvertTimeToDatetime64) Rebind(sql string, args ...driver.NamedValue) (string, []driver.NamedValue, error) {
	for i := range args {
		switch t := args[i].Value.(type) {
		case time.Time:
			args[i].Value = types.Datetime64Value(t.Unix())
		case *time.Time:
			args[i].Value = types.NullableDatetime64ValueFromTime(t)
		}
	}
	return sql, args, nil
}
