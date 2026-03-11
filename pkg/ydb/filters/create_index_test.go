package filters

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCreateIndex(t *testing.T) {
	for _, tt := range []struct {
		name string
		sql  string
		exp  string
	}{
		{
			name: "CREATE UNIQUE INDEX (backticks)",
			sql:  "CREATE UNIQUE INDEX `UQE_seed_assignment_builtin_role_role_name` ON `seed_assignment` (`builtin_role`, `role_name`);",
			exp:  "ALTER TABLE `seed_assignment` CREATE UNIQUE INDEX `UQE_seed_assignment_builtin_role_role_name` GLOBAL ON (`builtin_role`, `role_name`);",
		},
		{
			name: "CREATE UNIQUE INDEX (lower case)",
			sql:  "create unique index `UQE_seed_assignment_builtin_role_role_name` on `seed_assignment` (`builtin_role`, `role_name`);",
			exp:  "ALTER TABLE `seed_assignment` CREATE UNIQUE INDEX `UQE_seed_assignment_builtin_role_role_name` GLOBAL ON (`builtin_role`, `role_name`);",
		},
		{
			name: "CREATE UNIQUE INDEX (without quotas)",
			sql:  "CREATE UNIQUE INDEX UQE_seed_assignment_builtin_role_role_name ON seed_assignment (builtin_role, role_name);",
			exp:  "ALTER TABLE seed_assignment CREATE UNIQUE INDEX UQE_seed_assignment_builtin_role_role_name GLOBAL ON (builtin_role, role_name);",
		},
		{
			name: "CREATE INDEX (backticks)",
			sql:  "CREATE INDEX `UQE_seed_assignment_builtin_role_role_name` ON `seed_assignment` (`builtin_role`, `role_name`);",
			exp:  "ALTER TABLE `seed_assignment` CREATE INDEX `UQE_seed_assignment_builtin_role_role_name` GLOBAL ON (`builtin_role`, `role_name`);",
		},
		{
			name: "CREATE INDEX (without quotas)",
			sql:  "CREATE INDEX UQE_seed_assignment_builtin_role_role_name ON seed_assignment (builtin_role, role_name);",
			exp:  "ALTER TABLE seed_assignment CREATE INDEX UQE_seed_assignment_builtin_role_role_name GLOBAL ON (builtin_role, role_name);",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f := &CreateIndex{}
			require.Equal(t, tt.exp, f.Do(tt.sql, nil, nil))
		})
	}
}
