UPDATE annotation
	SET dashboard_uid = (SELECT uid FROM dashboard WHERE dashboard.id = annotation.dashboard_id)
	WHERE dashboard_uid IS NULL AND dashboard_id != 0 AND EXISTS (SELECT 1 FROM dashboard WHERE dashboard.id = annotation.dashboard_id);
