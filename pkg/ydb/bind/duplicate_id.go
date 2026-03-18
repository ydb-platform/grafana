package bind

import (
	"database/sql/driver"
	"strings"
)

var _ Binder = (*ReduceDuplicateIdInSelect)(nil)

type ReduceDuplicateIdInSelect struct{}

func (ReduceDuplicateIdInSelect) Rebind(sql string, args ...driver.NamedValue) (string, []driver.NamedValue, error) {
	return strings.ReplaceAll(sql, "SELECT `id`, `id`, ", "SELECT `id`, "), args, nil
}
