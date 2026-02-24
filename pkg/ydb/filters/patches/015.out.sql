UPSERT INTO anon_device (
  device_id,
  client_ip,
  user_agent,
  created_at,
  updated_at
)
VALUES
  (?, ?, ?, ?, ?)
