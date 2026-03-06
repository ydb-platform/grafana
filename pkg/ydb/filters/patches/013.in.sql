UPDATE dashboard_tag
	SET
    	dashboard_uid = (SELECT uid FROM dashboard WHERE dashboard.id = dashboard_tag.dashboard_id),
    	org_id = (SELECT org_id FROM dashboard WHERE dashboard.id = dashboard_tag.dashboard_id)
	WHERE
    	(dashboard_uid IS NULL OR org_id IS NULL)
    	AND EXISTS (SELECT 1 FROM dashboard WHERE dashboard.id = dashboard_tag.dashboard_id);
