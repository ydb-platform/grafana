package filters

import "github.com/grafana/grafana/pkg/util/xorm/core"

var (
	_ core.Filter = (*UpdateTableAlias)(nil)
)

type UpdateTableAlias struct{}

func (u UpdateTableAlias) Do(sql string, dialect core.Dialect, table *core.Table) string {
	//TODO implement me
	panic("implement me")
}
