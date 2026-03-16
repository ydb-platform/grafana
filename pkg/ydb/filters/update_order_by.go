package filters

import (
	"strings"

	"github.com/grafana/grafana/pkg/util/xorm/core"
)

var (
	_ core.Filter = (*RemoveOrderByFromUpdate)(nil)
)

type RemoveOrderByFromUpdate struct{}

func (RemoveOrderByFromUpdate) Do(sql string, dialect core.Dialect, table *core.Table) string {
	trimmed := strings.TrimSpace(sql)
	if len(trimmed) < 7 {
		return sql
	}
	up := trimmed[:7]
	if (up[0] != 'U' && up[0] != 'u') || (up[1] != 'P' && up[1] != 'p') || (up[2] != 'D' && up[2] != 'd') ||
		(up[3] != 'A' && up[3] != 'a') || (up[4] != 'T' && up[4] != 't') || (up[5] != 'E' && up[5] != 'e') || (up[6] != ' ' && up[6] != '\t') {
		return sql
	}
	idx := strings.Index(sql, " ORDER BY ")
	if idx < 0 {
		return sql
	}
	return sql[:idx]
}
