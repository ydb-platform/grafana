package filters

import (
	"database/sql"
	"testing"

	"github.com/grafana/grafana/pkg/util/xorm/core"
	"github.com/stretchr/testify/require"
)

func TestFilters(t *testing.T) {
	applyFilters := func(sql string, args []any) (string, []any) {
		for _, filter := range Filters {
			if filterWithArgs, has := filter.(core.FilterWithArgs); has {
				sql, args = filterWithArgs.DoWithArgs(sql, nil, &core.Table{
					PrimaryKeys: []string{"table_id"},
				}, args...)
			} else {
				sql = filter.Do(sql, nil, &core.Table{
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
			name: "base",
			in: args{
				sql:  "SELECT a, b, c FROM test_table WHERE p1=? AND (p2=? OR p3=?) AND (id) IN (?,?,?,?)",
				args: []any{1, 2, 3, 4, 5, 6, 7},
			},
			out: args{
				sql: "SELECT a, b, c FROM test_table WHERE p1=$p1 AND (p2=$p2 OR p3=$p3) AND table_id IN $p4",
				args: []any{
					sql.Named("p1", int64(1)),
					sql.Named("p2", int64(2)),
					sql.Named("p3", int64(3)),
					sql.Named("p4", []any{int64(4), int64(5), int64(6), int64(7)}),
				},
			},
		},
		{
			name: "dashboard",
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
					sql.Named("p1", int64(0)),
					sql.Named("p2", []any{
						int64(1), int64(2), int64(3), int64(4), int64(5), int64(6), int64(7), int64(8), int64(9), int64(0),
						int64(1), int64(2), int64(3), int64(4), int64(5), int64(6), int64(7), int64(8), int64(9), int64(0),
						int64(1), int64(2), int64(3), int64(4), int64(5), int64(6), int64(7), int64(8), int64(9), int64(0),
						int64(1), int64(2), int64(3), int64(4), int64(5), int64(6), int64(7), int64(8), int64(9), int64(0),
						int64(1), int64(2), int64(3), int64(4), int64(5), int64(6), int64(7), int64(8), int64(9), int64(0),
						int64(1), int64(2), int64(3), int64(4), int64(5), int64(6), int64(7), int64(8), int64(9),
					}),
					sql.Named("p3", int64(0)),
					sql.Named("p4", int64(1)),
				},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			sql, args := applyFilters(tt.in.sql, tt.in.args)
			require.Equal(t, tt.out.sql, sql)
			require.Equal(t, tt.out.args, args)
		})
	}
}
