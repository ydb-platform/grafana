package filters

import (
	"testing"

	"github.com/grafana/grafana/pkg/util/xorm/core"
	"github.com/stretchr/testify/require"
)

func TestID(t *testing.T) {
	for _, tt := range []struct {
		name  string
		sql   string
		table *core.Table
		exp   string
	}{
		{
			name: "",
			sql:  "SELECT a, b, c FROM test_table WHERE (id) IN (1,2,3)",
			table: &core.Table{
				PrimaryKeys: []string{"table_id"},
			},
			exp: "SELECT a, b, c FROM test_table WHERE table_id IN (1,2,3)",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f := &ReplaceIdPlaceholderToTableId{}
			require.Equal(t, tt.exp, f.Do(tt.sql, nil, tt.table))
		})
	}
}
