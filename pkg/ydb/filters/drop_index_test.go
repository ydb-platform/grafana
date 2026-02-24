package filters

import (
	"testing"

	"github.com/grafana/grafana/pkg/util/xorm/core"
	"github.com/stretchr/testify/require"
)

type dropIndexTestDialect struct {
	core.Dialect

	tables []*core.Table
}

func (d *dropIndexTestDialect) GetTables() ([]*core.Table, error) {
	return d.tables, nil
}

func TestDropIndex(t *testing.T) {
	for _, tt := range []struct {
		name    string
		sql     string
		dialect core.Dialect
		table   *core.Table
		exp     string
	}{
		{
			name:    "DROP INDEX (table not nil)",
			sql:     "DROP INDEX UQE_seed_assignment_builtin_role_role_name;",
			dialect: nil,
			table: &core.Table{
				Name: "seed_assignment",
			},
			exp: "ALTER TABLE `seed_assignment` DROP INDEX UQE_seed_assignment_builtin_role_role_name;",
		},
		{
			name:    "drop index (lower case)",
			sql:     "drop index UQE_seed_assignment_builtin_role_role_name;",
			dialect: nil,
			table: &core.Table{
				Name: "seed_assignment",
			},
			exp: "ALTER TABLE `seed_assignment` DROP INDEX UQE_seed_assignment_builtin_role_role_name;",
		},
		{
			name: "DROP INDEX (describe tables and test table names)",
			sql:  "DROP INDEX UQE_seed_assignment_builtin_role_role_name;",
			dialect: &dropIndexTestDialect{
				tables: []*core.Table{
					{Name: "seed_assignment"},
				},
			},
			table: nil,
			exp:   "ALTER TABLE `seed_assignment` DROP INDEX UQE_seed_assignment_builtin_role_role_name;",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f := &DropIndex{}
			require.Equal(t, tt.exp, f.Do(tt.sql, tt.dialect, tt.table))
		})
	}
}
