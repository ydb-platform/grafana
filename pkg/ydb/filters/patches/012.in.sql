UPDATE library_element
	SET folder_uid = dashboard.uid
	FROM dashboard
	WHERE library_element.folder_id = dashboard.id AND library_element.org_id = dashboard.org_id
