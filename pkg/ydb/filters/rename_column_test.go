package filters

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRenameColumn(t *testing.T) {
	for _, tt := range []struct {
		name string
		sql  string
		exp  string
	}{
		{
			name: "RENAME COLUMN (UPPER CASE)",
			sql:  "ALTER TABLE alert_instance RENAME COLUMN def_org_id TO rule_org_id;",
			exp:  "ALTER TABLE alert_instance ADD COLUMN rule_org_id Int64;\nUPDATE alert_instance SET rule_org_id = def_org_id;\nALTER TABLE alert_instance DROP COLUMN def_org_id;\n",
		},
		{
			name: "RENAME COLUMN (lower case)",
			sql:  "alter table alert_instance rename column def_org_id to rule_org_id;",
			exp:  "ALTER TABLE alert_instance ADD COLUMN rule_org_id Int64;\nUPDATE alert_instance SET rule_org_id = def_org_id;\nALTER TABLE alert_instance DROP COLUMN def_org_id;\n",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f := &RenameColumn{}
			require.Equal(t, tt.exp, f.Do(tt.sql, nil, nil))
		})
	}
}
