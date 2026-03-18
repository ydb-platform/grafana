package bind

import (
	"database/sql/driver"
	"strings"
)

var _ Binder = (*RemoveOrderByFromUpdate)(nil)

type RemoveOrderByFromUpdate struct{}

func (RemoveOrderByFromUpdate) Rebind(sql string, args ...driver.NamedValue) (string, []driver.NamedValue, error) {
	trimmed := strings.TrimSpace(sql)
	if len(trimmed) < 7 {
		return sql, args, nil
	}
	up := trimmed[:7]
	if (up[0] != 'U' && up[0] != 'u') || (up[1] != 'P' && up[1] != 'p') || (up[2] != 'D' && up[2] != 'd') ||
		(up[3] != 'A' && up[3] != 'a') || (up[4] != 'T' && up[4] != 't') || (up[5] != 'E' && up[5] != 'e') || (up[6] != ' ' && up[6] != '\t') {
		return sql, args, nil
	}
	idx := strings.Index(sql, " ORDER BY ")
	if idx < 0 {
		return sql, args, nil
	}
	return sql[:idx], args, nil
}
