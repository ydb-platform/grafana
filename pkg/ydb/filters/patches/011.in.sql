UPDATE preferences
	SET home_dashboard_uid = (SELECT uid FROM dashboard WHERE dashboard.id = preferences.home_dashboard_id)
	WHERE home_dashboard_uid IS NULL AND EXISTS (SELECT 1 FROM dashboard WHERE dashboard.id = preferences.home_dashboard_id);
