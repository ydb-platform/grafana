
	SELECT res.uid, res.is_folder, res.org_id
	FROM (SELECT dashboard.id, dashboard.uid, dashboard.is_folder, dashboard.org_id, count(dashboard_acl.id) as count
		  FROM dashboard
				LEFT JOIN dashboard_acl ON dashboard.id = dashboard_acl.dashboard_id
		  WHERE dashboard.has_acl IS TRUE
		  GROUP BY dashboard.id) as res
	WHERE res.count = 0

