package xorm

import (
	"database/sql"
	"reflect"
	"testing"

	"github.com/grafana/grafana/pkg/util/xorm/core"
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

func TestYdbInClauseFilter(t *testing.T) {
	type args struct {
		sql  string
		args []any
	}
	for _, tt := range []struct {
		name string
		in   args
		out  args
	}{
		{
			name: "IN keywork in literal (single quotas)",
			in: args{
				sql:  "SELECT ' IN (?,?,?)' WHERE ?, ?, ?",
				args: []any{1, 2, 3},
			},
			out: args{
				sql:  "SELECT ' IN (?,?,?)' WHERE ?, ?, ?",
				args: []any{1, 2, 3},
			},
		},
		{
			name: "IN keywork in literal (double quotas)",
			in: args{
				sql:  "SELECT \" IN (?,?,?)\" WHERE ?, ?, ?",
				args: []any{1, 2, 3},
			},
			out: args{
				sql:  "SELECT \" IN (?,?,?)\" WHERE ?, ?, ?",
				args: []any{1, 2, 3},
			},
		},
		{
			name: "single IN without separate spaces",
			in: args{
				sql:  "SELECT * FROM tbl WHERE id IN (?,?,?)",
				args: []any{1, 2, 3},
			},
			out: args{
				sql:  "SELECT * FROM tbl WHERE id IN ?",
				args: []any{[]any{1, 2, 3}},
			},
		},
		{
			name: "single IN with separate spaces",
			in: args{
				sql:  "SELECT * FROM tbl WHERE id IN (?, ?, ?)",
				args: []any{1, 2, 3},
			},
			out: args{
				sql:  "SELECT * FROM tbl WHERE id IN ?",
				args: []any{[]any{1, 2, 3}},
			},
		},
		{
			name: "single IN with single arg",
			in: args{
				sql:  "SELECT * FROM tbl WHERE id IN(?)",
				args: []any{1},
			},
			out: args{
				sql:  "SELECT * FROM tbl WHERE id IN ?",
				args: []any{[]any{1}},
			},
		},
		{
			name: "single IN with and without separate spaces",
			in: args{
				sql:  "SELECT * FROM tbl WHERE id IN(?,?, ?)",
				args: []any{1, 2, 3},
			},
			out: args{
				sql:  "SELECT * FROM tbl WHERE id IN ?",
				args: []any{[]any{1, 2, 3}},
			},
		},
		{
			name: "single IN between args",
			in: args{
				sql:  "SELECT * FROM tbl WHERE name = ? AND id IN(?,?, ?) OR deleted = ?",
				args: []any{"test", 1, 2, 3, "false"},
			},
			out: args{
				sql:  "SELECT * FROM tbl WHERE name = ? AND id IN ? OR deleted = ?",
				args: []any{"test", []any{1, 2, 3}, "false"},
			},
		},
		{
			name: "two IN's between args",
			in: args{
				sql:  "SELECT * FROM tbl WHERE name = ? AND id IN(?,?, ?) AND deleted = ? AND role_id IN (?,?,?) AND role_name LIKE ?",
				args: []any{"test", 1, 2, 3, "false", 4, 5, 6, "admin%"},
			},
			out: args{
				sql:  "SELECT * FROM tbl WHERE name = ? AND id IN(?,?, ?) AND deleted = ? AND role_id IN (?,?,?) AND role_name LIKE ?",
				args: []any{"test", []any{1, 2, 3}, "false", []any{4, 5, 6}, "admin%"},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			s := &ydbInClauseFilter{}
			sql, args := s.DoWithArgs(tt.in.sql, nil, nil, tt.in.args...)
			if sql != tt.out.sql {
				t.Errorf("Do() got = %q, want %q", sql, tt.out.sql)
			}
			if !reflect.DeepEqual(args, tt.out.args) {
				t.Errorf("Do() got1 = %v, want %v", args, tt.out.args)
			}
		})
	}
}

func TestYdbSeqFilter(t *testing.T) {
	type args struct {
		sql  string
		args []any
	}
	for _, tt := range []struct {
		name string
		in   args
		out  args
	}{
		{
			name: "happy path",
			in: args{
				sql:  "SELECT * WHERE p1=? AND (p2=? OR p3=?)",
				args: []any{1, 2, 3},
			},
			out: args{
				sql: "SELECT * WHERE p1=$p1 AND (p2=$p2 OR p3=$p3)",
				args: []any{
					sql.Named("p1", 1),
					sql.Named("p2", 2),
					sql.Named("p3", 3),
				},
			},
		},
		{
			name: "exclude single quotas",
			in: args{
				sql:  "SELECT ' IN (?,?,?)' WHERE p1=? AND (p2=? OR p3=?)",
				args: []any{1, 2, 3},
			},
			out: args{
				sql: "SELECT ' IN (?,?,?)' WHERE p1=$p1 AND (p2=$p2 OR p3=$p3)",
				args: []any{
					sql.Named("p1", 1),
					sql.Named("p2", 2),
					sql.Named("p3", 3),
				},
			},
		},
		{
			name: "exclude double quotas",
			in: args{
				sql:  "SELECT \" IN (?,?,?)\" WHERE p1=? AND (p2=? OR p3=?)",
				args: []any{1, "2", 3},
			},
			out: args{
				sql: "SELECT \" IN (?,?,?)\" WHERE p1=$p1 AND (p2=$p2 OR p3=$p3)",
				args: []any{
					sql.Named("p1", 1),
					sql.Named("p2", "2"),
					sql.Named("p3", 3),
				},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			s := &ydbSeqFilter{}
			sql, args := s.DoWithArgs(tt.in.sql, nil, nil, tt.in.args...)
			if sql != tt.out.sql {
				t.Errorf("Do() got = %q, want %q", sql, tt.out.sql)
			}
			if !reflect.DeepEqual(args, tt.out.args) {
				t.Errorf("Do() got1 = %v, want %v", args, tt.out.args)
			}
		})
	}
}

func TestYdbFilters(t *testing.T) {
	applyFilters := func(sql string, args []any) (string, []any) {
		for _, filter := range ydbFilters {
			if filterWithArgs, has := filter.(core.FilterWithArgs); has {
				sql, args = filterWithArgs.DoWithArgs(sql, &ydbDialect{}, &core.Table{
					PrimaryKeys: []string{"table_id"},
				}, args...)
			} else {
				sql = filter.Do(sql, &ydbDialect{}, &core.Table{
					PrimaryKeys: []string{"table_id"},
				})
			}
		}
		return sql, args
	}
	type args struct {
		sql  string
		args []any
	}
	for _, tt := range []struct {
		name string
		in   args
		out  args
	}{
		{
			name: "happy path",
			in: args{
				sql:  "SELECT a, b, c FROM test_table WHERE p1=? AND (p2=? OR p3=?) AND (id) IN (?,?,?,?)",
				args: []any{1, 2, 3, 4, 5, 6, 7},
			},
			out: args{
				sql: "SELECT a, b, c FROM test_table WHERE p1=$p1 AND (p2=$p2 OR p3=$p3) AND `table_id` IN $p4",
				args: []any{
					sql.Named("p1", 1),
					sql.Named("p2", 2),
					sql.Named("p3", 3),
					sql.Named("p4", []any{4, 5, 6, 7}),
				},
			},
		},
		{
			name: "",
			in: args{
				sql: `SELECT
						  dashboard.id,
						  dashboard.org_id,
						  dashboard.uid,
						  dashboard.title,
						  dashboard.slug,
						  dashboard_tag.term,
						  dashboard.is_folder,
						  dashboard.folder_id,
						  dashboard.deleted,
						  folder.uid AS folder_uid,
						  folder.title AS folder_slug,
						  folder.title AS folder_title,
						  dashboard.title AS title
						FROM (
						  SELECT dashboard.id AS id, dashboard.title AS title
						  FROM dashboard
						  WHERE dashboard.org_id = ?
							AND dashboard.uid IN (?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?,?)
							AND dashboard.is_folder = true
							AND dashboard.folder_uid IS NULL
							AND dashboard.uid != ?
							AND (dashboard.folder_uid != ? OR dashboard.folder_uid IS NULL)
						  ORDER BY title ASC
						  LIMIT 50 OFFSET 0
						) AS ids
						INNER JOIN dashboard ON ids.id = dashboard.id
						LEFT OUTER JOIN folder ON folder.uid = dashboard.folder_uid AND folder.org_id = dashboard.org_id
						LEFT OUTER JOIN dashboard_tag ON dashboard.id = dashboard_tag.dashboard_id
						ORDER BY dashboard.title ASC`,
				args: []any{
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
					0, 1, 2, 3, 4, 5, 6, 7, 8, 9,
					0, 1,
				},
			},
			out: args{
				sql: `SELECT
						  dashboard.id,
						  dashboard.org_id,
						  dashboard.uid,
						  dashboard.title,
						  dashboard.slug,
						  dashboard_tag.term,
						  dashboard.is_folder,
						  dashboard.folder_id,
						  dashboard.deleted,
						  folder.uid AS folder_uid,
						  folder.title AS folder_slug,
						  folder.title AS folder_title,
						  dashboard.title AS title
						FROM (
						  SELECT dashboard.id AS id, dashboard.title AS title
						  FROM dashboard
						  WHERE dashboard.org_id = $p1
							AND dashboard.uid IN $p2
							AND dashboard.is_folder = true
							AND dashboard.folder_uid IS NULL
							AND dashboard.uid != $p3
							AND (dashboard.folder_uid != $p4 OR dashboard.folder_uid IS NULL)
						  ORDER BY title ASC
						  LIMIT 50 OFFSET 0
						) AS ids
						INNER JOIN dashboard ON ids.id = dashboard.id
						LEFT OUTER JOIN folder ON folder.uid = dashboard.folder_uid AND folder.org_id = dashboard.org_id
						LEFT OUTER JOIN dashboard_tag ON dashboard.id = dashboard_tag.dashboard_id
						ORDER BY dashboard.title ASC`,
				args: []any{
					sql.Named("p1", 0),
					sql.Named("p2", []any{
						1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
						1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
						1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
						1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
						1, 2, 3, 4, 5, 6, 7, 8, 9, 0,
						1, 2, 3, 4, 5, 6, 7, 8, 9,
					}),
					sql.Named("p3", 0),
					sql.Named("p4", 1),
				},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			sql, args := applyFilters(tt.in.sql, tt.in.args)
			if sql != tt.out.sql {
				t.Errorf("Do() got = %q, want %q", sql, tt.out.sql)
			}
			if !reflect.DeepEqual(args, tt.out.args) {
				t.Errorf("Do() got1 = %v, want %v", args, tt.out.args)
			}
		})
	}
}
