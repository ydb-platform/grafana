UPDATE dashboard
	SET folder_uid = folder.uid
	FROM dashboard folder
	WHERE dashboard.folder_id = folder.id
	  AND dashboard.is_folder = ?
