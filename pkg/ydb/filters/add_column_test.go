package filters

import (
	"testing"

	"github.com/grafana/grafana/pkg/util/xorm/core"
	"github.com/stretchr/testify/require"
)

func TestAddColumn(t *testing.T) {
	for _, tt := range []struct {
		name  string
		sql   string
		table *core.Table
		exp   string
	}{
		{
			name: "",
			sql:  "ALTER TABLE `data_keys` ADD COLUMN `name` String",
			table: func() *core.Table {
				table := core.NewEmptyTable()
				table.AddColumn(&core.Column{
					Name:      "name",
					TableName: "data_keys",
				})
				return table
			}(),
			exp: "SELECT * FROM (SELECT 1) WHERE false; -- nop",
		},
		{
			name: "",
			sql:  "alter table `data_keys` add column `name` String",
			table: func() *core.Table {
				return core.NewEmptyTable()
			}(),
			exp: "alter table `data_keys` add column `name` String",
		},
		{
			name: "",
			sql:  "ALTER TABLE `data_keys` ADD COLUMN `name1` String",
			table: func() *core.Table {
				table := core.NewEmptyTable()
				table.AddColumn(&core.Column{
					Name:      "name",
					TableName: "data_keys",
				})
				return table
			}(),
			exp: "ALTER TABLE `data_keys` ADD COLUMN `name1` String",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f := &AddColumn{}
			require.Equal(t, tt.exp, f.Do(tt.sql, nil, tt.table))
		})
	}
}
