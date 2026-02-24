package db

import (
	"github.com/grafana/grafana/pkg/services/sqlstore/migrator"
	"github.com/grafana/grafana/pkg/util/xorm/core"
	"github.com/grafana/grafana/pkg/util/xorm/ydb"
)

func init() {
	core.RegisterDriver("ydb", &ydb.Driver{})
	core.RegisterDialect("ydb", func() core.Dialect { return &ydb.Dialect{} })
	migrator.Register("ydb", ydb.NewMigrator)
}
