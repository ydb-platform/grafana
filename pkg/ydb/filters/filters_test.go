package filters

import (
	"database/sql"
	"testing"

	"github.com/grafana/grafana/pkg/util/xorm/core"
	"github.com/stretchr/testify/require"
	"github.com/ydb-platform/ydb-go-sdk/v3/pkg/xstring"
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
			name: "IN clause with literal values",
			in: args{
				sql:  "SELECT * FROM tabl WHERE created > ? AND status in (?, ?) OR is_folder = 1",
				args: []any{1, 2, 3},
			},
			out: args{
				sql: "SELECT * FROM tabl WHERE created > $p1 AND status IN $p2 OR is_folder",
				args: []any{
					sql.Named("p1", int64(1)),
					sql.Named("p2", []any{int64(2), int64(3)}),
				},
			},
		},
		{
			name: "SetCreatedForOutstandingInvites",
			in: args{
				sql:  "UPDATE `temp_user` SET created = ?, updated = ? WHERE created = '0' AND status in (?,?)",
				args: []any{1772043198, 1772043198, "SignUpStarted", "InvitePending"},
			},
			out: args{
				sql: "UPDATE `temp_user` SET created = $p1, updated = $p2 WHERE created = $p3 AND status IN $p4",
				args: []any{
					sql.Named("p1", int64(1772043198)),
					sql.Named("p2", int64(1772043198)),
					sql.Named("p3", xstring.ToBytes("0")),
					sql.Named("p4", []any{
						xstring.ToBytes("SignUpStarted"),
						xstring.ToBytes("InvitePending"),
					}),
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
		{
			name: "patch for insert into dashbard_version",
			in: args{
				sql:  "INSERT INTO dashboard_version\n(\n\tdashboard_id,\n\tversion,\n\tparent_version,\n\trestored_from,\n\tcreated,\n\tcreated_by,\n\tmessage,\n\tdata\n)\nSELECT\n\tdashboard.id,\n\tdashboard.version,\n\tdashboard.version,\n\tdashboard.version,\n\tdashboard.updated,\n\tCOALESCE(dashboard.updated_by, -1),\n\t'',\n\tdashboard.data\nFROM dashboard;",
				args: nil,
			},
			out: args{
				sql: "INSERT INTO dashboard_version (\n  dashboard_id,\n  version,\n  parent_version,\n  restored_from,\n  created,\n  created_by,\n  message,\n  `data`\n)\nSELECT\n  dashboard.id AS dashboard_id,\n  dashboard.version AS version,\n  dashboard.version AS parent_version,\n  dashboard.version AS restored_from,\n  dashboard.updated AS created,\n  COALESCE(dashboard.updated_by, -1) AS created_by,\n  $p1 AS message,\n  dashboard.data AS `data`\nFROM dashboard;\n",
				args: []any{
					sql.Named("p1", []byte{}),
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
