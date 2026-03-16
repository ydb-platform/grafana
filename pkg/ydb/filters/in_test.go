package filters

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConvertInArgsToList(t *testing.T) {
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
			name: "IN keyword in literal (single quotas)",
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
			name: "IN keyword in literal (double quotas)",
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
			name: "IN in lower case",
			in: args{
				sql:  "SELECT * FROM tabl WHERE created > ? AND status in (?, ?)",
				args: []any{1, 2, 3},
			},
			out: args{
				sql:  "SELECT * FROM tabl WHERE created > ? AND status IN ?",
				args: []any{1, []any{2, 3}},
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
			name: "IN keyword in multi-line",
			in: args{
				sql: `SELECT * WHERE id IN (
					?,
					?,
					?
				)`,
				args: []any{1, 2, 3},
			},
			out: args{
				sql:  `SELECT * WHERE id IN ?`,
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
				sql:  "SELECT * FROM tbl WHERE name = ? AND id IN ? AND deleted = ? AND role_id IN ? AND role_name LIKE ?",
				args: []any{"test", []any{1, 2, 3}, "false", []any{4, 5, 6}, "admin%"},
			},
		},
		{
			name: "many IN's",
			in: args{
				sql: `
SELECT
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
  folder.title AS folder_title
FROM
  (
	SELECT
	  dashboard.id
	FROM
	  dashboard
	WHERE
	  dashboard.org_id = ?
	  AND dashboard.uid != ?
	  AND (
		dashboard.folder_uid != ?
		OR dashboard.folder_uid IS NULL
	  )
	  AND (
		(
		  dashboard.uid IN (
			SELECT
			  identifier
			FROM
			  permission
			WHERE
			  kind = 'dashboards'
			  AND attribute = 'uid'
			  AND role_id IN(
				SELECT
				  id
				FROM
				  role
				  INNER JOIN (
					SELECT
					  ur.role_id
					FROM
					  user_role AS ur
					WHERE
					  ur.user_id = ?
					  AND (
						ur.org_id = ?
						OR ur.org_id = ?
					  )
					UNION
					SELECT
					  br.role_id
					FROM
					  builtin_role AS br
					WHERE
					  br.role IN (?, ?)
					  AND (
						br.org_id = ?
						OR br.org_id = ?
					  )
				  ) as all_role ON role.id = all_role.role_id
			  )
			  AND action IN (?, ?, ?, ?)
		  )
		  AND NOT dashboard.is_folder
		)
		OR (
		  (
			dashboard.org_id = ?
			AND dashboard.folder_id IN (
			  SELECT
				d.id
			  FROM
				dashboard d
				INNER JOIN folder f1 ON d.uid = f1.uid
				AND d.org_id = f1.org_id
			  WHERE
				f1.org_id = ?
				AND f1.uid IN (
				  SELECT
					identifier
				  FROM
					permission
				  WHERE
					kind = 'folders'
					AND attribute = 'uid'
					AND role_id IN(
					  SELECT
						id
					  FROM
						role
						INNER JOIN (
						  SELECT
							ur.role_id
						  FROM
							user_role AS ur
						  WHERE
							ur.user_id = ?
							AND (
							  ur.org_id = ?
							  OR ur.org_id = ?
							)
						  UNION
						  SELECT
							br.role_id
						  FROM
							builtin_role AS br
						  WHERE
							br.role IN (?, ?)
							AND (
							  br.org_id = ?
							  OR br.org_id = ?
							)
						) as all_role ON role.id = all_role.role_id
					)
					AND action IN (?, ?, ?, ?)
				)
			)
		  )
		  OR (
			dashboard.org_id = ?
			AND dashboard.folder_id IN (
			  SELECT
				d.id
			  FROM
				dashboard d
				INNER JOIN folder f1 ON d.uid = f1.uid
				AND d.org_id = f1.org_id
				INNER JOIN folder f2 ON f1.parent_uid = f2.uid
				AND f1.org_id = f2.org_id
			  WHERE
				f2.org_id = ?
				AND f2.uid IN (
				  SELECT
					identifier
				  FROM
					permission
				  WHERE
					kind = 'folders'
					AND attribute = 'uid'
					AND role_id IN(
					  SELECT
						id
					  FROM
						role
						INNER JOIN (
						  SELECT
							ur.role_id
						  FROM
							user_role AS ur
						  WHERE
							ur.user_id = ?
							AND (
							  ur.org_id = ?
							  OR ur.org_id = ?
							)
						  UNION
						  SELECT
							br.role_id
						  FROM
							builtin_role AS br
						  WHERE
							br.role IN (?, ?)
							AND (
							  br.org_id = ?
							  OR br.org_id = ?
							)
						) as all_role ON role.id = all_role.role_id
					)
					AND action IN (?, ?, ?, ?)
				)
			)
		  )
		  OR (
			dashboard.org_id = ?
			AND dashboard.folder_id IN (
			  SELECT
				d.id
			  FROM
				dashboard d
				INNER JOIN folder f1 ON d.uid = f1.uid
				AND d.org_id = f1.org_id
				INNER JOIN folder f2 ON f1.parent_uid = f2.uid
				AND f1.org_id = f2.org_id
				INNER JOIN folder f3 ON f2.parent_uid = f3.uid
				AND f2.org_id = f3.org_id
			  WHERE
				f3.org_id = ?
				AND f3.uid IN (
				  SELECT
					identifier
				  FROM
					permission
				  WHERE
					kind = 'folders'
					AND attribute = 'uid'
					AND role_id IN(
					  SELECT
						id
					  FROM
						role
						INNER JOIN (
						  SELECT
							ur.role_id
						  FROM
							user_role AS ur
						  WHERE
							ur.user_id = ?
							AND (
							  ur.org_id = ?
							  OR ur.org_id = ?
							)
						  UNION
						  SELECT
							br.role_id
						  FROM
							builtin_role AS br
						  WHERE
							br.role IN (?, ?)
							AND (
							  br.org_id = ?
							  OR br.org_id = ?
							)
						) as all_role ON role.id = all_role.role_id
					)
					AND action IN (?, ?, ?, ?)
				)
			)
		  )
		  OR (
			dashboard.org_id = ?
			AND dashboard.folder_id IN (
			  SELECT
				d.id
			  FROM
				dashboard d
				INNER JOIN folder f1 ON d.uid = f1.uid
				AND d.org_id = f1.org_id
				INNER JOIN folder f2 ON f1.parent_uid = f2.uid
				AND f1.org_id = f2.org_id
				INNER JOIN folder f3 ON f2.parent_uid = f3.uid
				AND f2.org_id = f3.org_id
				INNER JOIN folder f4 ON f3.parent_uid = f4.uid
				AND f3.org_id = f4.org_id
			  WHERE
				f4.org_id = ?
				AND f4.uid IN (
				  SELECT
					identifier
				  FROM
					permission
				  WHERE
					kind = 'folders'
					AND attribute = 'uid'
					AND role_id IN(
					  SELECT
						id
					  FROM
						role
						INNER JOIN (
						  SELECT
							ur.role_id
						  FROM
							user_role AS ur
						  WHERE
							ur.user_id = ?
							AND (
							  ur.org_id = ?
							  OR ur.org_id = ?
							)
						  UNION
						  SELECT
							br.role_id
						  FROM
							builtin_role AS br
						  WHERE
							br.role IN (?, ?)
							AND (
							  br.org_id = ?
							  OR br.org_id = ?
							)
						) as all_role ON role.id = all_role.role_id
					)
					AND action IN (?, ?, ?, ?)
				)
			)
		  )
		  OR (
			dashboard.org_id = ?
			AND dashboard.folder_id IN (
			  SELECT
				d.id
			  FROM
				dashboard d
				INNER JOIN folder f1 ON d.uid = f1.uid
				AND d.org_id = f1.org_id
				INNER JOIN folder f2 ON f1.parent_uid = f2.uid
				AND f1.org_id = f2.org_id
				INNER JOIN folder f3 ON f2.parent_uid = f3.uid
				AND f2.org_id = f3.org_id
				INNER JOIN folder f4 ON f3.parent_uid = f4.uid
				AND f3.org_id = f4.org_id
				INNER JOIN folder f5 ON f4.parent_uid = f5.uid
				AND f4.org_id = f5.org_id
			  WHERE
				f5.org_id = ?
				AND f5.uid IN (
				  SELECT
					identifier
				  FROM
					permission
				  WHERE
					kind = 'folders'
					AND attribute = 'uid'
					AND role_id IN(
					  SELECT
						id
					  FROM
						role
						INNER JOIN (
						  SELECT
							ur.role_id
						  FROM
							user_role AS ur
						  WHERE
							ur.user_id = ?
							AND (
							  ur.org_id = ?
							  OR ur.org_id = ?
							)
						  UNION
						  SELECT
							br.role_id
						  FROM
							builtin_role AS br
						  WHERE
							br.role IN (?, ?)
							AND (
							  br.org_id = ?
							  OR br.org_id = ?
							)
						) as all_role ON role.id = all_role.role_id
					)
					AND action IN (?, ?, ?, ?)
				)
			)
		  )
		  OR (
			dashboard.org_id = ?
			AND dashboard.folder_id IN (
			  SELECT
				d.id
			  FROM
				dashboard d
				INNER JOIN folder f1 ON d.uid = f1.uid
				AND d.org_id = f1.org_id
				INNER JOIN folder f2 ON f1.parent_uid = f2.uid
				AND f1.org_id = f2.org_id
				INNER JOIN folder f3 ON f2.parent_uid = f3.uid
				AND f2.org_id = f3.org_id
				INNER JOIN folder f4 ON f3.parent_uid = f4.uid
				AND f3.org_id = f4.org_id
				INNER JOIN folder f5 ON f4.parent_uid = f5.uid
				AND f4.org_id = f5.org_id
				INNER JOIN folder f6 ON f5.parent_uid = f6.uid
				AND f5.org_id = f6.org_id
			  WHERE
				f6.org_id = ?
				AND f6.uid IN (
				  SELECT
					identifier
				  FROM
					permission
				  WHERE
					kind = 'folders'
					AND attribute = 'uid'
					AND role_id IN(
					  SELECT
						id
					  FROM
						role
						INNER JOIN (
						  SELECT
							ur.role_id
						  FROM
							user_role AS ur
						  WHERE
							ur.user_id = ?
							AND (
							  ur.org_id = ?
							  OR ur.org_id = ?
							)
						  UNION
						  SELECT
							br.role_id
						  FROM
							builtin_role AS br
						  WHERE
							br.role IN (?, ?)
							AND (
							  br.org_id = ?
							  OR br.org_id = ?
							)
						) as all_role ON role.id = all_role.role_id
					)
					AND action IN (?, ?, ?, ?)
				)
			)
		  )
		  AND NOT dashboard.is_folder
		)
		OR (
		  (
			dashboard.org_id = ?
			AND dashboard.uid IN (
			  SELECT
				d.uid
			  FROM
				dashboard d
				INNER JOIN folder f1 ON d.uid = f1.uid
				AND d.org_id = f1.org_id
			  WHERE
				f1.org_id = ?
				AND f1.uid IN (
				  SELECT
					identifier
				  FROM
					permission
				  WHERE
					kind = 'folders'
					AND attribute = 'uid'
					AND role_id IN(
					  SELECT
						id
					  FROM
						role
						INNER JOIN (
						  SELECT
							ur.role_id
						  FROM
							user_role AS ur
						  WHERE
							ur.user_id = ?
							AND (
							  ur.org_id = ?
							  OR ur.org_id = ?
							)
						  UNION
						  SELECT
							br.role_id
						  FROM
							builtin_role AS br
						  WHERE
							br.role IN (?, ?)
							AND (
							  br.org_id = ?
							  OR br.org_id = ?
							)
						) as all_role ON role.id = all_role.role_id
					)
					AND action IN (?, ?, ?, ?)
				)
			)
		  )
		  OR (
			dashboard.org_id = ?
			AND dashboard.uid IN (
			  SELECT
				d.uid
			  FROM
				dashboard d
				INNER JOIN folder f1 ON d.uid = f1.uid
				AND d.org_id = f1.org_id
				INNER JOIN folder f2 ON f1.parent_uid = f2.uid
				AND f1.org_id = f2.org_id
			  WHERE
				f2.org_id = ?
				AND f2.uid IN (
				  SELECT
					identifier
				  FROM
					permission
				  WHERE
					kind = 'folders'
					AND attribute = 'uid'
					AND role_id IN(
					  SELECT
						id
					  FROM
						role
						INNER JOIN (
						  SELECT
							ur.role_id
						  FROM
							user_role AS ur
						  WHERE
							ur.user_id = ?
							AND (
							  ur.org_id = ?
							  OR ur.org_id = ?
							)
						  UNION
						  SELECT
							br.role_id
						  FROM
							builtin_role AS br
						  WHERE
							br.role IN (?, ?)
							AND (
							  br.org_id = ?
							  OR br.org_id = ?
							)
						) as all_role ON role.id = all_role.role_id
					)
					AND action IN (?, ?, ?, ?)
				)
			)
		  )
		  OR (
			dashboard.org_id = ?
			AND dashboard.uid IN (
			  SELECT
				d.uid
			  FROM
				dashboard d
				INNER JOIN folder f1 ON d.uid = f1.uid
				AND d.org_id = f1.org_id
				INNER JOIN folder f2 ON f1.parent_uid = f2.uid
				AND f1.org_id = f2.org_id
				INNER JOIN folder f3 ON f2.parent_uid = f3.uid
				AND f2.org_id = f3.org_id
			  WHERE
				f3.org_id = ?
				AND f3.uid IN (
				  SELECT
					identifier
				  FROM
					permission
				  WHERE
					kind = 'folders'
					AND attribute = 'uid'
					AND role_id IN(
					  SELECT
						id
					  FROM
						role
						INNER JOIN (
						  SELECT
							ur.role_id
						  FROM
							user_role AS ur
						  WHERE
							ur.user_id = ?
							AND (
							  ur.org_id = ?
							  OR ur.org_id = ?
							)
						  UNION
						  SELECT
							br.role_id
						  FROM
							builtin_role AS br
						  WHERE
							br.role IN (?, ?)
							AND (
							  br.org_id = ?
							  OR br.org_id = ?
							)
						) as all_role ON role.id = all_role.role_id
					)
					AND action IN (?, ?, ?, ?)
				)
			)
		  )
		  OR (
			dashboard.org_id = ?
			AND dashboard.uid IN (
			  SELECT
				d.uid
			  FROM
				dashboard d
				INNER JOIN folder f1 ON d.uid = f1.uid
				AND d.org_id = f1.org_id
				INNER JOIN folder f2 ON f1.parent_uid = f2.uid
				AND f1.org_id = f2.org_id
				INNER JOIN folder f3 ON f2.parent_uid = f3.uid
				AND f2.org_id = f3.org_id
				INNER JOIN folder f4 ON f3.parent_uid = f4.uid
				AND f3.org_id = f4.org_id
			  WHERE
				f4.org_id = ?
				AND f4.uid IN (
				  SELECT
					identifier
				  FROM
					permission
				  WHERE
					kind = 'folders'
					AND attribute = 'uid'
					AND role_id IN(
					  SELECT
						id
					  FROM
						role
						INNER JOIN (
						  SELECT
							ur.role_id
						  FROM
							user_role AS ur
						  WHERE
							ur.user_id = ?
							AND (
							  ur.org_id = ?
							  OR ur.org_id = ?
							)
						  UNION
						  SELECT
							br.role_id
						  FROM
							builtin_role AS br
						  WHERE
							br.role IN (?, ?)
							AND (
							  br.org_id = ?
							  OR br.org_id = ?
							)
						) as all_role ON role.id = all_role.role_id
					)
					AND action IN (?, ?, ?, ?)
				)
			)
		  )
		  OR (
			dashboard.org_id = ?
			AND dashboard.uid IN (
			  SELECT
				d.uid
			  FROM
				dashboard d
				INNER JOIN folder f1 ON d.uid = f1.uid
				AND d.org_id = f1.org_id
				INNER JOIN folder f2 ON f1.parent_uid = f2.uid
				AND f1.org_id = f2.org_id
				INNER JOIN folder f3 ON f2.parent_uid = f3.uid
				AND f2.org_id = f3.org_id
				INNER JOIN folder f4 ON f3.parent_uid = f4.uid
				AND f3.org_id = f4.org_id
				INNER JOIN folder f5 ON f4.parent_uid = f5.uid
				AND f4.org_id = f5.org_id
			  WHERE
				f5.org_id = ?
				AND f5.uid IN (
				  SELECT
					identifier
				  FROM
					permission
				  WHERE
					kind = 'folders'
					AND attribute = 'uid'
					AND role_id IN(
					  SELECT
						id
					  FROM
						role
						INNER JOIN (
						  SELECT
							ur.role_id
						  FROM
							user_role AS ur
						  WHERE
							ur.user_id = ?
							AND (
							  ur.org_id = ?
							  OR ur.org_id = ?
							)
						  UNION
						  SELECT
							br.role_id
						  FROM
							builtin_role AS br
						  WHERE
							br.role IN (?, ?)
							AND (
							  br.org_id = ?
							  OR br.org_id = ?
							)
						) as all_role ON role.id = all_role.role_id
					)
					AND action IN (?, ?, ?, ?)
				)
			)
		  )
		  OR (
			dashboard.org_id = ?
			AND dashboard.uid IN (
			  SELECT
				d.uid
			  FROM
				dashboard d
				INNER JOIN folder f1 ON d.uid = f1.uid
				AND d.org_id = f1.org_id
				INNER JOIN folder f2 ON f1.parent_uid = f2.uid
				AND f1.org_id = f2.org_id
				INNER JOIN folder f3 ON f2.parent_uid = f3.uid
				AND f2.org_id = f3.org_id
				INNER JOIN folder f4 ON f3.parent_uid = f4.uid
				AND f3.org_id = f4.org_id
				INNER JOIN folder f5 ON f4.parent_uid = f5.uid
				AND f4.org_id = f5.org_id
				INNER JOIN folder f6 ON f5.parent_uid = f6.uid
				AND f5.org_id = f6.org_id
			  WHERE
				f6.org_id = ?
				AND f6.uid IN (
				  SELECT
					identifier
				  FROM
					permission
				  WHERE
					kind = 'folders'
					AND attribute = 'uid'
					AND role_id IN(
					  SELECT
						id
					  FROM
						role
						INNER JOIN (
						  SELECT
							ur.role_id
						  FROM
							user_role AS ur
						  WHERE
							ur.user_id = ?
							AND (
							  ur.org_id = ?
							  OR ur.org_id = ?
							)
						  UNION
						  SELECT
							br.role_id
						  FROM
							builtin_role AS br
						  WHERE
							br.role IN (?, ?)
							AND (
							  br.org_id = ?
							  OR br.org_id = ?
							)
						) as all_role ON role.id = all_role.role_id
					)
					AND action IN (?, ?, ?, ?)
				)
			)
		  )
		  AND dashboard.is_folder
		)
	  )
	  AND dashboard.deleted IS NULL
	ORDER BY
	  dashboard.title ASC
	LIMIT
	  30 OFFSET 0
  ) AS ids
  INNER JOIN dashboard ON ids.id = dashboard.id
  LEFT OUTER JOIN folder ON folder.uid = dashboard.folder_uid
  AND folder.org_id = dashboard.org_id
  LEFT OUTER JOIN dashboard_tag ON dashboard.id = dashboard_tag.dashboard_id
ORDER BY
  dashboard.title ASC
				`,
				args: []any{
					"-1", "k6-app", "k6-app", "1", "-1", "0", "None", "Grafana Admin", "-1", "0", "dashboards:read", "dashboards:view", "dashboards:edit", "dashboards:admin", "-1", "-1", "1", "-1", "0", "None", "Grafana Admin", "-1", "0", "dashboards:read", "folders:view", "folders:edit", "folders:admin", "-1", "-1", "1", "-1", "0", "None", "Grafana Admin", "-1", "0", "dashboards:read", "folders:view", "folders:edit", "folders:admin", "-1", "-1", "1", "-1", "0", "None", "Grafana Admin", "-1", "0", "dashboards:read", "folders:view", "folders:edit", "folders:admin", "-1", "-1", "1", "-1", "0", "None", "Grafana Admin", "-1", "0", "dashboards:read", "folders:view", "folders:edit", "folders:admin", "-1", "-1", "1", "-1", "0", "None", "Grafana Admin", "-1", "0", "dashboards:read", "folders:view", "folders:edit", "folders:admin", "-1", "-1", "1", "-1", "0", "None", "Grafana Admin", "-1", "0", "dashboards:read", "folders:view", "folders:edit", "folders:admin", "-1", "-1", "1", "-1", "0", "None", "Grafana Admin", "-1", "0", "folders:read", "folders:view", "folders:edit", "folders:admin", "-1", "-1", "1", "-1", "0", "None", "Grafana Admin", "-1", "0", "folders:read", "folders:view", "folders:edit", "folders:admin", "-1", "-1", "1", "-1", "0", "None", "Grafana Admin", "-1", "0", "folders:read", "folders:view", "folders:edit", "folders:admin", "-1", "-1", "1", "-1", "0", "None", "Grafana Admin", "-1", "0", "folders:read", "folders:view", "folders:edit", "folders:admin", "-1", "-1", "1", "-1", "0", "None", "Grafana Admin", "-1", "0", "folders:read", "folders:view", "folders:edit", "folders:admin", "-1", "-1", "1", "-1", "0", "None", "Grafana Admin", "-1", "0", "folders:read", "folders:view", "folders:edit", "folders:admin",
				},
			},
			out: args{
				sql: `
SELECT
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
  folder.title AS folder_title
FROM
  (
	SELECT
	  dashboard.id
	FROM
	  dashboard
	WHERE
	  dashboard.org_id = ?
	  AND dashboard.uid != ?
	  AND (
		dashboard.folder_uid != ?
		OR dashboard.folder_uid IS NULL
	  )
	  AND (
		(
		  dashboard.uid IN (
			SELECT
			  identifier
			FROM
			  permission
			WHERE
			  kind = 'dashboards'
			  AND attribute = 'uid'
			  AND role_id IN(
				SELECT
				  id
				FROM
				  role
				  INNER JOIN (
					SELECT
					  ur.role_id
					FROM
					  user_role AS ur
					WHERE
					  ur.user_id = ?
					  AND (
						ur.org_id = ?
						OR ur.org_id = ?
					  )
					UNION
					SELECT
					  br.role_id
					FROM
					  builtin_role AS br
					WHERE
					  br.role IN ?
					  AND (
						br.org_id = ?
						OR br.org_id = ?
					  )
				  ) as all_role ON role.id = all_role.role_id
			  )
			  AND action IN ?
		  )
		  AND NOT dashboard.is_folder
		)
		OR (
		  (
			dashboard.org_id = ?
			AND dashboard.folder_id IN (
			  SELECT
				d.id
			  FROM
				dashboard d
				INNER JOIN folder f1 ON d.uid = f1.uid
				AND d.org_id = f1.org_id
			  WHERE
				f1.org_id = ?
				AND f1.uid IN (
				  SELECT
					identifier
				  FROM
					permission
				  WHERE
					kind = 'folders'
					AND attribute = 'uid'
					AND role_id IN(
					  SELECT
						id
					  FROM
						role
						INNER JOIN (
						  SELECT
							ur.role_id
						  FROM
							user_role AS ur
						  WHERE
							ur.user_id = ?
							AND (
							  ur.org_id = ?
							  OR ur.org_id = ?
							)
						  UNION
						  SELECT
							br.role_id
						  FROM
							builtin_role AS br
						  WHERE
							br.role IN ?
							AND (
							  br.org_id = ?
							  OR br.org_id = ?
							)
						) as all_role ON role.id = all_role.role_id
					)
					AND action IN ?
				)
			)
		  )
		  OR (
			dashboard.org_id = ?
			AND dashboard.folder_id IN (
			  SELECT
				d.id
			  FROM
				dashboard d
				INNER JOIN folder f1 ON d.uid = f1.uid
				AND d.org_id = f1.org_id
				INNER JOIN folder f2 ON f1.parent_uid = f2.uid
				AND f1.org_id = f2.org_id
			  WHERE
				f2.org_id = ?
				AND f2.uid IN (
				  SELECT
					identifier
				  FROM
					permission
				  WHERE
					kind = 'folders'
					AND attribute = 'uid'
					AND role_id IN(
					  SELECT
						id
					  FROM
						role
						INNER JOIN (
						  SELECT
							ur.role_id
						  FROM
							user_role AS ur
						  WHERE
							ur.user_id = ?
							AND (
							  ur.org_id = ?
							  OR ur.org_id = ?
							)
						  UNION
						  SELECT
							br.role_id
						  FROM
							builtin_role AS br
						  WHERE
							br.role IN ?
							AND (
							  br.org_id = ?
							  OR br.org_id = ?
							)
						) as all_role ON role.id = all_role.role_id
					)
					AND action IN ?
				)
			)
		  )
		  OR (
			dashboard.org_id = ?
			AND dashboard.folder_id IN (
			  SELECT
				d.id
			  FROM
				dashboard d
				INNER JOIN folder f1 ON d.uid = f1.uid
				AND d.org_id = f1.org_id
				INNER JOIN folder f2 ON f1.parent_uid = f2.uid
				AND f1.org_id = f2.org_id
				INNER JOIN folder f3 ON f2.parent_uid = f3.uid
				AND f2.org_id = f3.org_id
			  WHERE
				f3.org_id = ?
				AND f3.uid IN (
				  SELECT
					identifier
				  FROM
					permission
				  WHERE
					kind = 'folders'
					AND attribute = 'uid'
					AND role_id IN(
					  SELECT
						id
					  FROM
						role
						INNER JOIN (
						  SELECT
							ur.role_id
						  FROM
							user_role AS ur
						  WHERE
							ur.user_id = ?
							AND (
							  ur.org_id = ?
							  OR ur.org_id = ?
							)
						  UNION
						  SELECT
							br.role_id
						  FROM
							builtin_role AS br
						  WHERE
							br.role IN ?
							AND (
							  br.org_id = ?
							  OR br.org_id = ?
							)
						) as all_role ON role.id = all_role.role_id
					)
					AND action IN ?
				)
			)
		  )
		  OR (
			dashboard.org_id = ?
			AND dashboard.folder_id IN (
			  SELECT
				d.id
			  FROM
				dashboard d
				INNER JOIN folder f1 ON d.uid = f1.uid
				AND d.org_id = f1.org_id
				INNER JOIN folder f2 ON f1.parent_uid = f2.uid
				AND f1.org_id = f2.org_id
				INNER JOIN folder f3 ON f2.parent_uid = f3.uid
				AND f2.org_id = f3.org_id
				INNER JOIN folder f4 ON f3.parent_uid = f4.uid
				AND f3.org_id = f4.org_id
			  WHERE
				f4.org_id = ?
				AND f4.uid IN (
				  SELECT
					identifier
				  FROM
					permission
				  WHERE
					kind = 'folders'
					AND attribute = 'uid'
					AND role_id IN(
					  SELECT
						id
					  FROM
						role
						INNER JOIN (
						  SELECT
							ur.role_id
						  FROM
							user_role AS ur
						  WHERE
							ur.user_id = ?
							AND (
							  ur.org_id = ?
							  OR ur.org_id = ?
							)
						  UNION
						  SELECT
							br.role_id
						  FROM
							builtin_role AS br
						  WHERE
							br.role IN ?
							AND (
							  br.org_id = ?
							  OR br.org_id = ?
							)
						) as all_role ON role.id = all_role.role_id
					)
					AND action IN ?
				)
			)
		  )
		  OR (
			dashboard.org_id = ?
			AND dashboard.folder_id IN (
			  SELECT
				d.id
			  FROM
				dashboard d
				INNER JOIN folder f1 ON d.uid = f1.uid
				AND d.org_id = f1.org_id
				INNER JOIN folder f2 ON f1.parent_uid = f2.uid
				AND f1.org_id = f2.org_id
				INNER JOIN folder f3 ON f2.parent_uid = f3.uid
				AND f2.org_id = f3.org_id
				INNER JOIN folder f4 ON f3.parent_uid = f4.uid
				AND f3.org_id = f4.org_id
				INNER JOIN folder f5 ON f4.parent_uid = f5.uid
				AND f4.org_id = f5.org_id
			  WHERE
				f5.org_id = ?
				AND f5.uid IN (
				  SELECT
					identifier
				  FROM
					permission
				  WHERE
					kind = 'folders'
					AND attribute = 'uid'
					AND role_id IN(
					  SELECT
						id
					  FROM
						role
						INNER JOIN (
						  SELECT
							ur.role_id
						  FROM
							user_role AS ur
						  WHERE
							ur.user_id = ?
							AND (
							  ur.org_id = ?
							  OR ur.org_id = ?
							)
						  UNION
						  SELECT
							br.role_id
						  FROM
							builtin_role AS br
						  WHERE
							br.role IN ?
							AND (
							  br.org_id = ?
							  OR br.org_id = ?
							)
						) as all_role ON role.id = all_role.role_id
					)
					AND action IN ?
				)
			)
		  )
		  OR (
			dashboard.org_id = ?
			AND dashboard.folder_id IN (
			  SELECT
				d.id
			  FROM
				dashboard d
				INNER JOIN folder f1 ON d.uid = f1.uid
				AND d.org_id = f1.org_id
				INNER JOIN folder f2 ON f1.parent_uid = f2.uid
				AND f1.org_id = f2.org_id
				INNER JOIN folder f3 ON f2.parent_uid = f3.uid
				AND f2.org_id = f3.org_id
				INNER JOIN folder f4 ON f3.parent_uid = f4.uid
				AND f3.org_id = f4.org_id
				INNER JOIN folder f5 ON f4.parent_uid = f5.uid
				AND f4.org_id = f5.org_id
				INNER JOIN folder f6 ON f5.parent_uid = f6.uid
				AND f5.org_id = f6.org_id
			  WHERE
				f6.org_id = ?
				AND f6.uid IN (
				  SELECT
					identifier
				  FROM
					permission
				  WHERE
					kind = 'folders'
					AND attribute = 'uid'
					AND role_id IN(
					  SELECT
						id
					  FROM
						role
						INNER JOIN (
						  SELECT
							ur.role_id
						  FROM
							user_role AS ur
						  WHERE
							ur.user_id = ?
							AND (
							  ur.org_id = ?
							  OR ur.org_id = ?
							)
						  UNION
						  SELECT
							br.role_id
						  FROM
							builtin_role AS br
						  WHERE
							br.role IN ?
							AND (
							  br.org_id = ?
							  OR br.org_id = ?
							)
						) as all_role ON role.id = all_role.role_id
					)
					AND action IN ?
				)
			)
		  )
		  AND NOT dashboard.is_folder
		)
		OR (
		  (
			dashboard.org_id = ?
			AND dashboard.uid IN (
			  SELECT
				d.uid
			  FROM
				dashboard d
				INNER JOIN folder f1 ON d.uid = f1.uid
				AND d.org_id = f1.org_id
			  WHERE
				f1.org_id = ?
				AND f1.uid IN (
				  SELECT
					identifier
				  FROM
					permission
				  WHERE
					kind = 'folders'
					AND attribute = 'uid'
					AND role_id IN(
					  SELECT
						id
					  FROM
						role
						INNER JOIN (
						  SELECT
							ur.role_id
						  FROM
							user_role AS ur
						  WHERE
							ur.user_id = ?
							AND (
							  ur.org_id = ?
							  OR ur.org_id = ?
							)
						  UNION
						  SELECT
							br.role_id
						  FROM
							builtin_role AS br
						  WHERE
							br.role IN ?
							AND (
							  br.org_id = ?
							  OR br.org_id = ?
							)
						) as all_role ON role.id = all_role.role_id
					)
					AND action IN ?
				)
			)
		  )
		  OR (
			dashboard.org_id = ?
			AND dashboard.uid IN (
			  SELECT
				d.uid
			  FROM
				dashboard d
				INNER JOIN folder f1 ON d.uid = f1.uid
				AND d.org_id = f1.org_id
				INNER JOIN folder f2 ON f1.parent_uid = f2.uid
				AND f1.org_id = f2.org_id
			  WHERE
				f2.org_id = ?
				AND f2.uid IN (
				  SELECT
					identifier
				  FROM
					permission
				  WHERE
					kind = 'folders'
					AND attribute = 'uid'
					AND role_id IN(
					  SELECT
						id
					  FROM
						role
						INNER JOIN (
						  SELECT
							ur.role_id
						  FROM
							user_role AS ur
						  WHERE
							ur.user_id = ?
							AND (
							  ur.org_id = ?
							  OR ur.org_id = ?
							)
						  UNION
						  SELECT
							br.role_id
						  FROM
							builtin_role AS br
						  WHERE
							br.role IN ?
							AND (
							  br.org_id = ?
							  OR br.org_id = ?
							)
						) as all_role ON role.id = all_role.role_id
					)
					AND action IN ?
				)
			)
		  )
		  OR (
			dashboard.org_id = ?
			AND dashboard.uid IN (
			  SELECT
				d.uid
			  FROM
				dashboard d
				INNER JOIN folder f1 ON d.uid = f1.uid
				AND d.org_id = f1.org_id
				INNER JOIN folder f2 ON f1.parent_uid = f2.uid
				AND f1.org_id = f2.org_id
				INNER JOIN folder f3 ON f2.parent_uid = f3.uid
				AND f2.org_id = f3.org_id
			  WHERE
				f3.org_id = ?
				AND f3.uid IN (
				  SELECT
					identifier
				  FROM
					permission
				  WHERE
					kind = 'folders'
					AND attribute = 'uid'
					AND role_id IN(
					  SELECT
						id
					  FROM
						role
						INNER JOIN (
						  SELECT
							ur.role_id
						  FROM
							user_role AS ur
						  WHERE
							ur.user_id = ?
							AND (
							  ur.org_id = ?
							  OR ur.org_id = ?
							)
						  UNION
						  SELECT
							br.role_id
						  FROM
							builtin_role AS br
						  WHERE
							br.role IN ?
							AND (
							  br.org_id = ?
							  OR br.org_id = ?
							)
						) as all_role ON role.id = all_role.role_id
					)
					AND action IN ?
				)
			)
		  )
		  OR (
			dashboard.org_id = ?
			AND dashboard.uid IN (
			  SELECT
				d.uid
			  FROM
				dashboard d
				INNER JOIN folder f1 ON d.uid = f1.uid
				AND d.org_id = f1.org_id
				INNER JOIN folder f2 ON f1.parent_uid = f2.uid
				AND f1.org_id = f2.org_id
				INNER JOIN folder f3 ON f2.parent_uid = f3.uid
				AND f2.org_id = f3.org_id
				INNER JOIN folder f4 ON f3.parent_uid = f4.uid
				AND f3.org_id = f4.org_id
			  WHERE
				f4.org_id = ?
				AND f4.uid IN (
				  SELECT
					identifier
				  FROM
					permission
				  WHERE
					kind = 'folders'
					AND attribute = 'uid'
					AND role_id IN(
					  SELECT
						id
					  FROM
						role
						INNER JOIN (
						  SELECT
							ur.role_id
						  FROM
							user_role AS ur
						  WHERE
							ur.user_id = ?
							AND (
							  ur.org_id = ?
							  OR ur.org_id = ?
							)
						  UNION
						  SELECT
							br.role_id
						  FROM
							builtin_role AS br
						  WHERE
							br.role IN ?
							AND (
							  br.org_id = ?
							  OR br.org_id = ?
							)
						) as all_role ON role.id = all_role.role_id
					)
					AND action IN ?
				)
			)
		  )
		  OR (
			dashboard.org_id = ?
			AND dashboard.uid IN (
			  SELECT
				d.uid
			  FROM
				dashboard d
				INNER JOIN folder f1 ON d.uid = f1.uid
				AND d.org_id = f1.org_id
				INNER JOIN folder f2 ON f1.parent_uid = f2.uid
				AND f1.org_id = f2.org_id
				INNER JOIN folder f3 ON f2.parent_uid = f3.uid
				AND f2.org_id = f3.org_id
				INNER JOIN folder f4 ON f3.parent_uid = f4.uid
				AND f3.org_id = f4.org_id
				INNER JOIN folder f5 ON f4.parent_uid = f5.uid
				AND f4.org_id = f5.org_id
			  WHERE
				f5.org_id = ?
				AND f5.uid IN (
				  SELECT
					identifier
				  FROM
					permission
				  WHERE
					kind = 'folders'
					AND attribute = 'uid'
					AND role_id IN(
					  SELECT
						id
					  FROM
						role
						INNER JOIN (
						  SELECT
							ur.role_id
						  FROM
							user_role AS ur
						  WHERE
							ur.user_id = ?
							AND (
							  ur.org_id = ?
							  OR ur.org_id = ?
							)
						  UNION
						  SELECT
							br.role_id
						  FROM
							builtin_role AS br
						  WHERE
							br.role IN ?
							AND (
							  br.org_id = ?
							  OR br.org_id = ?
							)
						) as all_role ON role.id = all_role.role_id
					)
					AND action IN ?
				)
			)
		  )
		  OR (
			dashboard.org_id = ?
			AND dashboard.uid IN (
			  SELECT
				d.uid
			  FROM
				dashboard d
				INNER JOIN folder f1 ON d.uid = f1.uid
				AND d.org_id = f1.org_id
				INNER JOIN folder f2 ON f1.parent_uid = f2.uid
				AND f1.org_id = f2.org_id
				INNER JOIN folder f3 ON f2.parent_uid = f3.uid
				AND f2.org_id = f3.org_id
				INNER JOIN folder f4 ON f3.parent_uid = f4.uid
				AND f3.org_id = f4.org_id
				INNER JOIN folder f5 ON f4.parent_uid = f5.uid
				AND f4.org_id = f5.org_id
				INNER JOIN folder f6 ON f5.parent_uid = f6.uid
				AND f5.org_id = f6.org_id
			  WHERE
				f6.org_id = ?
				AND f6.uid IN (
				  SELECT
					identifier
				  FROM
					permission
				  WHERE
					kind = 'folders'
					AND attribute = 'uid'
					AND role_id IN(
					  SELECT
						id
					  FROM
						role
						INNER JOIN (
						  SELECT
							ur.role_id
						  FROM
							user_role AS ur
						  WHERE
							ur.user_id = ?
							AND (
							  ur.org_id = ?
							  OR ur.org_id = ?
							)
						  UNION
						  SELECT
							br.role_id
						  FROM
							builtin_role AS br
						  WHERE
							br.role IN ?
							AND (
							  br.org_id = ?
							  OR br.org_id = ?
							)
						) as all_role ON role.id = all_role.role_id
					)
					AND action IN ?
				)
			)
		  )
		  AND dashboard.is_folder
		)
	  )
	  AND dashboard.deleted IS NULL
	ORDER BY
	  dashboard.title ASC
	LIMIT
	  30 OFFSET 0
  ) AS ids
  INNER JOIN dashboard ON ids.id = dashboard.id
  LEFT OUTER JOIN folder ON folder.uid = dashboard.folder_uid
  AND folder.org_id = dashboard.org_id
  LEFT OUTER JOIN dashboard_tag ON dashboard.id = dashboard_tag.dashboard_id
ORDER BY
  dashboard.title ASC
				`,
				args: []any{
					"-1", "k6-app", "k6-app", "1", "-1", "0",
					[]any{"None", "Grafana Admin"}, "-1", "0", []any{"dashboards:read", "dashboards:view", "dashboards:edit", "dashboards:admin"}, "-1", "-1", "1", "-1", "0",
					[]any{"None", "Grafana Admin"}, "-1", "0", []any{"dashboards:read", "folders:view", "folders:edit", "folders:admin"}, "-1", "-1", "1", "-1", "0",
					[]any{"None", "Grafana Admin"}, "-1", "0", []any{"dashboards:read", "folders:view", "folders:edit", "folders:admin"}, "-1", "-1", "1", "-1", "0",
					[]any{"None", "Grafana Admin"}, "-1", "0", []any{"dashboards:read", "folders:view", "folders:edit", "folders:admin"}, "-1", "-1", "1", "-1", "0",
					[]any{"None", "Grafana Admin"}, "-1", "0", []any{"dashboards:read", "folders:view", "folders:edit", "folders:admin"}, "-1", "-1", "1", "-1", "0",
					[]any{"None", "Grafana Admin"}, "-1", "0", []any{"dashboards:read", "folders:view", "folders:edit", "folders:admin"}, "-1", "-1", "1", "-1", "0",
					[]any{"None", "Grafana Admin"}, "-1", "0", []any{"dashboards:read", "folders:view", "folders:edit", "folders:admin"}, "-1", "-1", "1", "-1", "0",
					[]any{"None", "Grafana Admin"}, "-1", "0", []any{"folders:read", "folders:view", "folders:edit", "folders:admin"}, "-1", "-1", "1", "-1", "0",
					[]any{"None", "Grafana Admin"}, "-1", "0", []any{"folders:read", "folders:view", "folders:edit", "folders:admin"}, "-1", "-1", "1", "-1", "0",
					[]any{"None", "Grafana Admin"}, "-1", "0", []any{"folders:read", "folders:view", "folders:edit", "folders:admin"}, "-1", "-1", "1", "-1", "0",
					[]any{"None", "Grafana Admin"}, "-1", "0", []any{"folders:read", "folders:view", "folders:edit", "folders:admin"}, "-1", "-1", "1", "-1", "0",
					[]any{"None", "Grafana Admin"}, "-1", "0", []any{"folders:read", "folders:view", "folders:edit", "folders:admin"}, "-1", "-1", "1", "-1", "0",
					[]any{"None", "Grafana Admin"}, "-1", "0", []any{"folders:read", "folders:view", "folders:edit", "folders:admin"},
				},
			},
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			f := &ConvertInArgsToList{}
			sql, args := f.DoWithArgs(tt.in.sql, nil, nil, tt.in.args...)
			require.Equal(t, tt.out.sql, sql)
			require.Equal(t, tt.out.args, args)
		})
	}
}
