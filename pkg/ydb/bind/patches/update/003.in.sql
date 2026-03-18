UPDATE star
	SET
    	dashboard_uid = (SELECT uid FROM dashboard WHERE dashboard.id = star.dashboard_id),
    	org_id = (SELECT org_id FROM dashboard WHERE dashboard.id = star.dashboard_id),
    	updated = DATETIME('now')
	WHERE
    	(dashboard_uid IS NULL OR org_id IS NULL)
    	AND EXISTS (SELECT 1 FROM dashboard WHERE dashboard.id = star.dashboard_id);
