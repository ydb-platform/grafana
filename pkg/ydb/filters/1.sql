PRAGMA ydb.CostBasedOptimization = "on";
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
		 FROM ( SELECT dashboard.id FROM dashboard  WHERE dashboard.org_id=? AND dashboard.uid != ? AND (dashboard.folder_uid != ? OR dashboard.folder_uid IS NULL) AND ((dashboard.uid IN (SELECT identifier FROM permission WHERE kind = 'dashboards' AND attribute = 'uid' AND role_id IN(SELECT id FROM role INNER JOIN (
			SELECT ur.role_id
			FROM user_role AS ur
			WHERE ur.user_id = ?
			AND (ur.org_id = ? OR ur.org_id = ?)
		UNION
			SELECT br.role_id FROM builtin_role AS br
			WHERE br.role IN (?, ?)
			AND (br.org_id = ? OR br.org_id = ?)
		) as all_role ON role.id = all_role.role_id)  AND action IN (?, ?, ?, ?)) AND NOT dashboard.is_folder) OR ((dashboard.org_id = ? AND dashboard.folder_id IN (SELECT d.id FROM dashboard d INNER JOIN folder f1 ON d.uid = f1.uid AND d.org_id = f1.org_id  WHERE f1.org_id = ? AND f1.uid IN (SELECT identifier FROM permission WHERE kind = 'folders' AND attribute = 'uid' AND role_id IN(SELECT id FROM role INNER JOIN (
			SELECT ur.role_id
			FROM user_role AS ur
			WHERE ur.user_id = ?
			AND (ur.org_id = ? OR ur.org_id = ?)
		UNION
			SELECT br.role_id FROM builtin_role AS br
			WHERE br.role IN (?, ?)
			AND (br.org_id = ? OR br.org_id = ?)
		) as all_role ON role.id = all_role.role_id)  AND action IN (?, ?, ?, ?)))) OR (dashboard.org_id = ? AND dashboard.folder_id IN (SELECT d.id FROM dashboard d INNER JOIN folder f1 ON d.uid = f1.uid AND d.org_id = f1.org_id  INNER JOIN folder f2 ON f1.parent_uid = f2.uid AND f1.org_id = f2.org_id  WHERE f2.org_id = ? AND f2.uid IN (SELECT identifier FROM permission WHERE kind = 'folders' AND attribute = 'uid' AND role_id IN(SELECT id FROM role INNER JOIN (
			SELECT ur.role_id
			FROM user_role AS ur
			WHERE ur.user_id = ?
			AND (ur.org_id = ? OR ur.org_id = ?)
		UNION
			SELECT br.role_id FROM builtin_role AS br
			WHERE br.role IN (?, ?)
			AND (br.org_id = ? OR br.org_id = ?)
		) as all_role ON role.id = all_role.role_id)  AND action IN (?, ?, ?, ?)))) OR (dashboard.org_id = ? AND dashboard.folder_id IN (SELECT d.id FROM dashboard d INNER JOIN folder f1 ON d.uid = f1.uid AND d.org_id = f1.org_id  INNER JOIN folder f2 ON f1.parent_uid = f2.uid AND f1.org_id = f2.org_id  INNER JOIN folder f3 ON f2.parent_uid = f3.uid AND f2.org_id = f3.org_id  WHERE f3.org_id = ? AND f3.uid IN (SELECT identifier FROM permission WHERE kind = 'folders' AND attribute = 'uid' AND role_id IN(SELECT id FROM role INNER JOIN (
			SELECT ur.role_id
			FROM user_role AS ur
			WHERE ur.user_id = ?
			AND (ur.org_id = ? OR ur.org_id = ?)
		UNION
			SELECT br.role_id FROM builtin_role AS br
			WHERE br.role IN (?, ?)
			AND (br.org_id = ? OR br.org_id = ?)
		) as all_role ON role.id = all_role.role_id)  AND action IN (?, ?, ?, ?)))) OR (dashboard.org_id = ? AND dashboard.folder_id IN (SELECT d.id FROM dashboard d INNER JOIN folder f1 ON d.uid = f1.uid AND d.org_id = f1.org_id  INNER JOIN folder f2 ON f1.parent_uid = f2.uid AND f1.org_id = f2.org_id  INNER JOIN folder f3 ON f2.parent_uid = f3.uid AND f2.org_id = f3.org_id  INNER JOIN folder f4 ON f3.parent_uid = f4.uid AND f3.org_id = f4.org_id  WHERE f4.org_id = ? AND f4.uid IN (SELECT identifier FROM permission WHERE kind = 'folders' AND attribute = 'uid' AND role_id IN(SELECT id FROM role INNER JOIN (
			SELECT ur.role_id
			FROM user_role AS ur
			WHERE ur.user_id = ?
			AND (ur.org_id = ? OR ur.org_id = ?)
		UNION
			SELECT br.role_id FROM builtin_role AS br
			WHERE br.role IN (?, ?)
			AND (br.org_id = ? OR br.org_id = ?)
		) as all_role ON role.id = all_role.role_id)  AND action IN (?, ?, ?, ?)))) OR (dashboard.org_id = ? AND dashboard.folder_id IN (SELECT d.id FROM dashboard d INNER JOIN folder f1 ON d.uid = f1.uid AND d.org_id = f1.org_id  INNER JOIN folder f2 ON f1.parent_uid = f2.uid AND f1.org_id = f2.org_id  INNER JOIN folder f3 ON f2.parent_uid = f3.uid AND f2.org_id = f3.org_id  INNER JOIN folder f4 ON f3.parent_uid = f4.uid AND f3.org_id = f4.org_id  INNER JOIN folder f5 ON f4.parent_uid = f5.uid AND f4.org_id = f5.org_id  WHERE f5.org_id = ? AND f5.uid IN (SELECT identifier FROM permission WHERE kind = 'folders' AND attribute = 'uid' AND role_id IN(SELECT id FROM role INNER JOIN (
			SELECT ur.role_id
			FROM user_role AS ur
			WHERE ur.user_id = ?
			AND (ur.org_id = ? OR ur.org_id = ?)
		UNION
			SELECT br.role_id FROM builtin_role AS br
			WHERE br.role IN (?, ?)
			AND (br.org_id = ? OR br.org_id = ?)
		) as all_role ON role.id = all_role.role_id)  AND action IN (?, ?, ?, ?)))) OR (dashboard.org_id = ? AND dashboard.folder_id IN (SELECT d.id FROM dashboard d INNER JOIN folder f1 ON d.uid = f1.uid AND d.org_id = f1.org_id  INNER JOIN folder f2 ON f1.parent_uid = f2.uid AND f1.org_id = f2.org_id  INNER JOIN folder f3 ON f2.parent_uid = f3.uid AND f2.org_id = f3.org_id  INNER JOIN folder f4 ON f3.parent_uid = f4.uid AND f3.org_id = f4.org_id  INNER JOIN folder f5 ON f4.parent_uid = f5.uid AND f4.org_id = f5.org_id  INNER JOIN folder f6 ON f5.parent_uid = f6.uid AND f5.org_id = f6.org_id  WHERE f6.org_id = ? AND f6.uid IN (SELECT identifier FROM permission WHERE kind = 'folders' AND attribute = 'uid' AND role_id IN(SELECT id FROM role INNER JOIN (
			SELECT ur.role_id
			FROM user_role AS ur
			WHERE ur.user_id = ?
			AND (ur.org_id = ? OR ur.org_id = ?)
		UNION
			SELECT br.role_id FROM builtin_role AS br
			WHERE br.role IN (?, ?)
			AND (br.org_id = ? OR br.org_id = ?)
		) as all_role ON role.id = all_role.role_id)  AND action IN (?, ?, ?, ?)))) AND NOT dashboard.is_folder) OR ((dashboard.org_id = ? AND dashboard.uid IN (SELECT d.uid FROM dashboard d INNER JOIN folder f1 ON d.uid = f1.uid AND d.org_id = f1.org_id  WHERE f1.org_id = ? AND f1.uid IN (SELECT identifier FROM permission WHERE kind = 'folders' AND attribute = 'uid' AND role_id IN(SELECT id FROM role INNER JOIN (
			SELECT ur.role_id
			FROM user_role AS ur
			WHERE ur.user_id = ?
			AND (ur.org_id = ? OR ur.org_id = ?)
		UNION
			SELECT br.role_id FROM builtin_role AS br
			WHERE br.role IN (?, ?)
			AND (br.org_id = ? OR br.org_id = ?)
		) as all_role ON role.id = all_role.role_id)  AND action IN (?, ?, ?, ?)))) OR (dashboard.org_id = ? AND dashboard.uid IN (SELECT d.uid FROM dashboard d INNER JOIN folder f1 ON d.uid = f1.uid AND d.org_id = f1.org_id  INNER JOIN folder f2 ON f1.parent_uid = f2.uid AND f1.org_id = f2.org_id  WHERE f2.org_id = ? AND f2.uid IN (SELECT identifier FROM permission WHERE kind = 'folders' AND attribute = 'uid' AND role_id IN(SELECT id FROM role INNER JOIN (
			SELECT ur.role_id
			FROM user_role AS ur
			WHERE ur.user_id = ?
			AND (ur.org_id = ? OR ur.org_id = ?)
		UNION
			SELECT br.role_id FROM builtin_role AS br
			WHERE br.role IN (?, ?)
			AND (br.org_id = ? OR br.org_id = ?)
		) as all_role ON role.id = all_role.role_id)  AND action IN (?, ?, ?, ?)))) OR (dashboard.org_id = ? AND dashboard.uid IN (SELECT d.uid FROM dashboard d INNER JOIN folder f1 ON d.uid = f1.uid AND d.org_id = f1.org_id  INNER JOIN folder f2 ON f1.parent_uid = f2.uid AND f1.org_id = f2.org_id  INNER JOIN folder f3 ON f2.parent_uid = f3.uid AND f2.org_id = f3.org_id  WHERE f3.org_id = ? AND f3.uid IN (SELECT identifier FROM permission WHERE kind = 'folders' AND attribute = 'uid' AND role_id IN(SELECT id FROM role INNER JOIN (
			SELECT ur.role_id
			FROM user_role AS ur
			WHERE ur.user_id = ?
			AND (ur.org_id = ? OR ur.org_id = ?)
		UNION
			SELECT br.role_id FROM builtin_role AS br
			WHERE br.role IN (?, ?)
			AND (br.org_id = ? OR br.org_id = ?)
		) as all_role ON role.id = all_role.role_id)  AND action IN (?, ?, ?, ?)))) OR (dashboard.org_id = ? AND dashboard.uid IN (SELECT d.uid FROM dashboard d INNER JOIN folder f1 ON d.uid = f1.uid AND d.org_id = f1.org_id  INNER JOIN folder f2 ON f1.parent_uid = f2.uid AND f1.org_id = f2.org_id  INNER JOIN folder f3 ON f2.parent_uid = f3.uid AND f2.org_id = f3.org_id  INNER JOIN folder f4 ON f3.parent_uid = f4.uid AND f3.org_id = f4.org_id  WHERE f4.org_id = ? AND f4.uid IN (SELECT identifier FROM permission WHERE kind = 'folders' AND attribute = 'uid' AND role_id IN(SELECT id FROM role INNER JOIN (
			SELECT ur.role_id
			FROM user_role AS ur
			WHERE ur.user_id = ?
			AND (ur.org_id = ? OR ur.org_id = ?)
		UNION
			SELECT br.role_id FROM builtin_role AS br
			WHERE br.role IN (?, ?)
			AND (br.org_id = ? OR br.org_id = ?)
		) as all_role ON role.id = all_role.role_id)  AND action IN (?, ?, ?, ?)))) OR (dashboard.org_id = ? AND dashboard.uid IN (SELECT d.uid FROM dashboard d INNER JOIN folder f1 ON d.uid = f1.uid AND d.org_id = f1.org_id  INNER JOIN folder f2 ON f1.parent_uid = f2.uid AND f1.org_id = f2.org_id  INNER JOIN folder f3 ON f2.parent_uid = f3.uid AND f2.org_id = f3.org_id  INNER JOIN folder f4 ON f3.parent_uid = f4.uid AND f3.org_id = f4.org_id  INNER JOIN folder f5 ON f4.parent_uid = f5.uid AND f4.org_id = f5.org_id  WHERE f5.org_id = ? AND f5.uid IN (SELECT identifier FROM permission WHERE kind = 'folders' AND attribute = 'uid' AND role_id IN(SELECT id FROM role INNER JOIN (
			SELECT ur.role_id
			FROM user_role AS ur
			WHERE ur.user_id = ?
			AND (ur.org_id = ? OR ur.org_id = ?)
		UNION
			SELECT br.role_id FROM builtin_role AS br
			WHERE br.role IN (?, ?)
			AND (br.org_id = ? OR br.org_id = ?)
		) as all_role ON role.id = all_role.role_id)  AND action IN (?, ?, ?, ?)))) OR (dashboard.org_id = ? AND dashboard.uid IN (SELECT d.uid FROM dashboard d INNER JOIN folder f1 ON d.uid = f1.uid AND d.org_id = f1.org_id  INNER JOIN folder f2 ON f1.parent_uid = f2.uid AND f1.org_id = f2.org_id  INNER JOIN folder f3 ON f2.parent_uid = f3.uid AND f2.org_id = f3.org_id  INNER JOIN folder f4 ON f3.parent_uid = f4.uid AND f3.org_id = f4.org_id  INNER JOIN folder f5 ON f4.parent_uid = f5.uid AND f4.org_id = f5.org_id  INNER JOIN folder f6 ON f5.parent_uid = f6.uid AND f5.org_id = f6.org_id  WHERE f6.org_id = ? AND f6.uid IN (SELECT identifier FROM permission WHERE kind = 'folders' AND attribute = 'uid' AND role_id IN(SELECT id FROM role INNER JOIN (
			SELECT ur.role_id
			FROM user_role AS ur
			WHERE ur.user_id = ?
			AND (ur.org_id = ? OR ur.org_id = ?)
		UNION
			SELECT br.role_id FROM builtin_role AS br
			WHERE br.role IN (?, ?)
			AND (br.org_id = ? OR br.org_id = ?)
		) as all_role ON role.id = all_role.role_id)  AND action IN (?, ?, ?, ?)))) AND dashboard.is_folder)) AND dashboard.deleted IS NULL ORDER BY dashboard.title ASC LIMIT 30 OFFSET 0) AS ids
		INNER JOIN dashboard ON ids.id = dashboard.id
LEFT OUTER JOIN folder ON folder.uid = dashboard.folder_uid AND folder.org_id = dashboard.org_id
	LEFT OUTER JOIN dashboard_tag ON dashboard.id = dashboard_tag.dashboard_id
 ORDER BY dashboard.title ASC
