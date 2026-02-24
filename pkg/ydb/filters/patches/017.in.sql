			DELETE FROM folder WHERE NOT EXISTS
				(SELECT 1 FROM dashboard WHERE dashboard.uid = folder.uid AND dashboard.org_id = folder.org_id AND dashboard.is_folder = true)
		
