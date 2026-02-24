package db

import (
	"github.com/grafana/grafana/pkg/services/sqlstore/migrator"
	"github.com/grafana/grafana/pkg/util/xorm/core"
	ydb2 "github.com/grafana/grafana/pkg/ydb"
)

func init() {
	core.RegisterDriver("ydb", &ydb2.Driver{})
	core.RegisterDialect("ydb", func() core.Dialect { return &ydb2.Dialect{} })
	migrator.Register("ydb", ydb2.NewMigrator)
}
