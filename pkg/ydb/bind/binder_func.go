package bind

import "database/sql/driver"

var _ Binder = Func(nil)

type Func func(sql string, args ...driver.NamedValue) (string, []driver.NamedValue, error)

func (f Func) Rebind(sql string, args ...driver.NamedValue) (string, []driver.NamedValue, error) {
	return f(sql, args...)
}
