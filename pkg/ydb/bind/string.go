package bind

import (
	"database/sql/driver"
	"time"

	"github.com/ydb-platform/ydb-go-sdk/v3/table/types"
)

var _ Binder = (*ConvertStringToDatetime64)(nil)

type ConvertStringToDatetime64 struct{}

func (f *ConvertStringToDatetime64) Rebind(sql string, args ...driver.NamedValue) (string, []driver.NamedValue, error) {
	for i := range args {
		switch v := args[i].Value.(type) {
		case string:
			if t, err := time.Parse(time.DateTime, v); err == nil {
				args[i].Value = types.Datetime64Value(t.Unix())
			}
		}
	}
	return sql, args, nil
}
