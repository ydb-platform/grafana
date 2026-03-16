package filters

import (
	"strings"

	"github.com/grafana/grafana/pkg/util/xorm/core"
)

// ReplaceIdPlaceholderToTableId filter SQL replace (id) to primary key column name
type ReplaceIdPlaceholderToTableId struct{}

func (i *ReplaceIdPlaceholderToTableId) Do(sql string, dialect core.Dialect, table *core.Table) string {
	if table != nil && len(table.PrimaryKeys) == 1 {
		return strings.ReplaceAll(sql, " (id) ", " "+table.PrimaryKeys[0]+" ")
	}
	return sql
}
