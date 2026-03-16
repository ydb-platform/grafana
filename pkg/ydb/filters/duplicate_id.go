package filters

import (
	"strings"

	"github.com/grafana/grafana/pkg/util/xorm/core"
)

var (
	_ core.Filter = (*ReduceDuplicateIdInSelect)(nil)
)

type ReduceDuplicateIdInSelect struct{}

func (ReduceDuplicateIdInSelect) Do(sql string, dialect core.Dialect, table *core.Table) string {
	return strings.ReplaceAll(sql, "SELECT `id`, `id`, ", "SELECT `id`, ")
}
