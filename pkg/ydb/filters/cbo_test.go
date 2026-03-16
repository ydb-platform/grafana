package filters

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestShouldUseCostBasedOptimization(t *testing.T) {
	tests := []struct {
		name                           string
		query                          string
		shouldUseCostBasedOptimization bool
	}{
		{
			"simple select",
			"SELECT * FROM cache_data WHERE key = ?",
			false,
		},
		{
			"server_lock",
			"SELECT * FROM server_lock WHERE operation_uid = ?",
			false,
		},
		{
			"dashboard list with permission",
			"SELECT id FROM dashboard WHERE org_id=? AND dashboard.uid IN (SELECT identifier FROM permission WHERE kind = 'dashboards') ORDER BY title LIMIT 50",
			true,
		},
		{
			"folder and permission",
			"SELECT id FROM dashboard WHERE folder_id IN (SELECT id FROM folder WHERE uid IN (SELECT identifier FROM permission WHERE kind = 'folders')) ORDER BY title",
			true,
		},
		{
			"two IN (SELECT",
			"SELECT a FROM t WHERE x IN (SELECT 1) AND y IN (SELECT 2 FROM u)",
			true,
		},
		{
			"one IN (SELECT",
			"SELECT * FROM permission WHERE role_id IN (SELECT id FROM role)",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cbo := EnableCostBasedOptimizer{}
			require.Equal(t, tt.shouldUseCostBasedOptimization, cbo.shouldUseCostBasedOptimization(tt.query))
		})
	}
}

func TestCostBasedOptimizer_Do(t *testing.T) {
	for _, tt := range []struct {
		name string
		in   string
		out  string
	}{
		{
			"simple select",
			"SELECT * FROM cache_data WHERE key = ?",
			"SELECT * FROM cache_data WHERE key = ?",
		},
		{
			"server_lock",
			"SELECT * FROM server_lock WHERE operation_uid = ?",
			"SELECT * FROM server_lock WHERE operation_uid = ?",
		},
		{
			"dashboard list with permission",
			"SELECT id FROM dashboard WHERE org_id=? AND dashboard.uid IN (SELECT identifier FROM permission WHERE kind = 'dashboards') ORDER BY title LIMIT 50",
			"PRAGMA ydb.CostBasedOptimization = \"on\";\nSELECT id FROM dashboard WHERE org_id=? AND dashboard.uid IN (SELECT identifier FROM permission WHERE kind = 'dashboards') ORDER BY title LIMIT 50",
		},
		{
			"folder and permission",
			"SELECT id FROM dashboard WHERE folder_id IN (SELECT id FROM folder WHERE uid IN (SELECT identifier FROM permission WHERE kind = 'folders')) ORDER BY title",
			"PRAGMA ydb.CostBasedOptimization = \"on\";\nSELECT id FROM dashboard WHERE folder_id IN (SELECT id FROM folder WHERE uid IN (SELECT identifier FROM permission WHERE kind = 'folders')) ORDER BY title",
		},
		{
			"two IN (SELECT",
			"SELECT a FROM t WHERE x IN (SELECT 1) AND y IN (SELECT 2 FROM u)",
			"PRAGMA ydb.CostBasedOptimization = \"on\";\nSELECT a FROM t WHERE x IN (SELECT 1) AND y IN (SELECT 2 FROM u)",
		},
		{
			"one IN (SELECT",
			"SELECT * FROM permission WHERE role_id IN (SELECT id FROM role)",
			"SELECT * FROM permission WHERE role_id IN (SELECT id FROM role)",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f := &EnableCostBasedOptimizer{}
			require.Equal(t, tt.out, f.Do(tt.in, nil, nil))
		})
	}
}
