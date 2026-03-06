package main

import (
	"github.com/grafana/grafana/pkg/services/sqlstore/migrator"
	"github.com/grafana/grafana/pkg/util/xorm/core"
	"github.com/grafana/grafana/pkg/ydb"
)

func init() {
	ydb := ydb.NewYDB()

	core.RegisterDriver("ydb", ydb)
	core.RegisterDialect("ydb", func() core.Dialect {
		return ydb
	})
	migrator.RegisterDialect("ydb", func() migrator.Dialect {
		return ydb
	})
}
