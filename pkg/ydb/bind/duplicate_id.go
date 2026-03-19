package bind

import "strings"

var _ Binder = (*ReduceDuplicateIdInSelect)(nil)

type ReduceDuplicateIdInSelect struct{}

func (ReduceDuplicateIdInSelect) Rebind(sql string, args ...any) (string, []any, error) {
	return strings.ReplaceAll(sql, "SELECT `id`, `id`, ", "SELECT `id`, "), args, nil
}
