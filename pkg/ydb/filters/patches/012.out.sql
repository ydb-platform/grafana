INSERT INTO
  server_lock (operation_uid, last_execution, version)
VALUES
  (?, ?, ?) RETURNING id
