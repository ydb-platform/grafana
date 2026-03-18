package bind

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUpdateOrderBy(t *testing.T) {
	for _, tt := range []struct {
		name string
		sql  string
		exp  string
	}{
		{
			name: "update where order by (parentheses)",
			sql:  "UPDATE `alert_configuration_history` SET `last_applied` = ? WHERE (org_id = ? AND configuration_hash = ?) ORDER BY `id`",
			exp:  "UPDATE `alert_configuration_history` SET `last_applied` = ? WHERE (org_id = ? AND configuration_hash = ?)",
		},
		{
			name: "update where order by (no parentheses)",
			sql:  "UPDATE `alert_configuration_history` SET `last_applied` = ? WHERE org_id = ? AND configuration_hash = ? ORDER BY `id`",
			exp:  "UPDATE `alert_configuration_history` SET `last_applied` = ? WHERE org_id = ? AND configuration_hash = ?",
		},
		{
			name: "update where order by desc",
			sql:  "UPDATE `alert_configuration_history` SET `last_applied` = ? WHERE (org_id = ? AND configuration_hash = ?) ORDER BY `id` DESC",
			exp:  "UPDATE `alert_configuration_history` SET `last_applied` = ? WHERE (org_id = ? AND configuration_hash = ?)",
		},
		{
			name: "update where order by asc",
			sql:  "UPDATE `alert_configuration_history` SET `last_applied` = ? WHERE (org_id = ? AND configuration_hash = ?) ORDER BY `id` ASC",
			exp:  "UPDATE `alert_configuration_history` SET `last_applied` = ? WHERE (org_id = ? AND configuration_hash = ?)",
		},
		{
			name: "no op for select (parentheses)",
			sql:  "SELECT * FROM `alert_configuration_history` SET `last_applied` = ? WHERE (org_id = ? AND configuration_hash = ?) ORDER BY `id` DESC",
			exp:  "SELECT * FROM `alert_configuration_history` SET `last_applied` = ? WHERE (org_id = ? AND configuration_hash = ?) ORDER BY `id` DESC",
		},
		{
			name: "no op for select (no parentheses)",
			sql:  "SELECT * FROM `alert_configuration_history` SET `last_applied` = ? WHERE org_id = ? AND configuration_hash = ? ORDER BY `id` DESC",
			exp:  "SELECT * FROM `alert_configuration_history` SET `last_applied` = ? WHERE org_id = ? AND configuration_hash = ? ORDER BY `id` DESC",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f := &RemoveOrderByFromUpdate{}
			out, _, err := f.Rebind(tt.sql)
			require.NoError(t, err)
			require.Equal(t, tt.exp, out)
		})
	}
}
