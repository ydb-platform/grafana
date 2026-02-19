package xorm

import (
	"testing"
)

func TestGetLastPartOfColumn(t *testing.T) {
	tests := []struct {
		column   string
		expected string
	}{
		{"table.column", "column"},
		{"column", "column"},
		{".", ""},
		{".x", "x"},
		{"xxxx.xxx.", ""},
	}

	for _, test := range tests {
		if getLastPartOfColumn(test.column) != test.expected {
			t.Errorf("Expected %s, got %s", test.expected, getLastPartOfColumn(test.column))
		}
	}
}

func TestRewriteILIKEToLowerLike(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{"no ilike", "SELECT * FROM t WHERE id = $1", "SELECT * FROM t WHERE id = $1"},
		{"dashboard title ilike", "WHERE dashboard.org_id=$1 AND dashboard.title ILIKE $p2", "WHERE dashboard.org_id=$1 AND LOWER(dashboard.title) LIKE LOWER($p2)"},
		{"table_col ilike", "AND le.name ILIKE $p3", "AND LOWER(le.name) LIKE LOWER($p3)"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := rewriteILIKEToLowerLike(tt.in)
			if got != tt.want {
				t.Errorf("rewriteILIKEToLowerLike() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestShouldUseCostBasedOptimization(t *testing.T) {
	tests := []struct {
		name  string
		query string
		want  bool
	}{
		{"simple select", "SELECT * FROM cache_data WHERE key = $1", false},
		{"server_lock", "SELECT * FROM server_lock WHERE operation_uid = $1", false},
		{"dashboard list with permission", "SELECT id FROM dashboard WHERE org_id=$1 AND dashboard.uid IN (SELECT identifier FROM permission WHERE kind = 'dashboards') ORDER BY title LIMIT 50", true},
		{"folder and permission", "SELECT id FROM dashboard WHERE folder_id IN (SELECT id FROM folder WHERE uid IN (SELECT identifier FROM permission WHERE kind = 'folders')) ORDER BY title", true},
		{"two IN (SELECT", "SELECT a FROM t WHERE x IN (SELECT 1) AND y IN (SELECT 2 FROM u)", true},
		{"one IN (SELECT", "SELECT * FROM permission WHERE role_id IN (SELECT id FROM role)", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := shouldUseCostBasedOptimization(tt.query)
			if got != tt.want {
				t.Errorf("shouldUseCostBasedOptimization() = %v, want %v", got, tt.want)
			}
		})
	}
}
