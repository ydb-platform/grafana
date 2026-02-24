INSERT INTO anon_device (device_id, client_ip, user_agent, created_at, updated_at)
					VALUES (?, ?, ?, ?, ?)
					ON CONFLICT (device_id) DO UPDATE SET
					client_ip = excluded.client_ip,
					user_agent = excluded.user_agent,
					updated_at = excluded.updated_at
