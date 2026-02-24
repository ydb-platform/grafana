package main

import (
	"github.com/grafana/grafana/pkg/services/sqlstore/migrator"
	"github.com/grafana/grafana/pkg/util/xorm/core"
	"github.com/grafana/grafana/pkg/ydb"
)

func init() {
	migrator.Register("ydb", ydb.NewMigrator)
	core.RegisterDriver("ydb", &ydb.Driver{})
	core.RegisterDialect("ydb", func() core.Dialect {
		return &ydb.Dialect{}
	})
}
