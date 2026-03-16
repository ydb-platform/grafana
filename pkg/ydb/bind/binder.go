package bind

import "database/sql/driver"

type Binder interface {
	Rebind(sql string, args ...driver.NamedValue) (string, []driver.NamedValue, error)
}
