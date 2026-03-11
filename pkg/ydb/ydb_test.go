package ydb

import (
	"testing"

	"github.com/grafana/grafana/pkg/services/sqlstore/migrator"
	"github.com/stretchr/testify/require"
)

func TestCreateIndexSQL(t *testing.T) {
	for _, tt := range []struct {
		name      string
		tableName string
		index     *migrator.Index
		sql       string
	}{
		{
			name:      "secondary index",
			tableName: "cache_data",
			index: &migrator.Index{
				Name: "cache_key",
				Type: migrator.UniqueIndex,
				Cols: []string{
					"cache_key",
				},
			},
			sql: "ALTER TABLE `cache_data` ADD INDEX `cache_data_cache_key_idx` GLOBAL ON (`cache_key`)",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			d := &YDB{}
			require.Equal(t, tt.sql, d.CreateIndexSQL(tt.tableName, tt.index))
		})
	}
}
