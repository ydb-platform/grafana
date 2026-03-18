package bind

import "database/sql/driver"

var _ Binder = Prepend("")

type Prepend string

func (p Prepend) Rebind(sql string, args ...driver.NamedValue) (string, []driver.NamedValue, error) {
	return string(p) + ";\n\n" + sql, args, nil
}
