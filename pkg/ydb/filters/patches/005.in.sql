UPDATE dashboard
	SET folder_uid = folder.parent_uid
	FROM folder
	WHERE dashboard.uid = folder.uid AND dashboard.org_id = folder.org_id
	AND dashboard.is_folder = ?
