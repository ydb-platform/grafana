INSERT INTO dashboard_version (
  dashboard_id,
  version,
  parent_version,
  restored_from,
  created,
  created_by,
  message,
  `data`
)
SELECT
  dashboard.id AS dashboard_id,
  dashboard.version AS version,
  dashboard.version AS parent_version,
  dashboard.version AS restored_from,
  dashboard.updated AS created,
  COALESCE(dashboard.updated_by, -1) AS created_by,
  '' AS message,
  dashboard.data AS `data`
FROM dashboard;
