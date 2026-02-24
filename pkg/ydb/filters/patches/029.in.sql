INSERT INTO dashboard_version
(
	dashboard_id,
	version,
	parent_version,
	restored_from,
	created,
	created_by,
	message,
	data
)
SELECT
	dashboard.id,
	dashboard.version,
	dashboard.version,
	dashboard.version,
	dashboard.updated,
	COALESCE(dashboard.updated_by, -1),
	'',
	dashboard.data
FROM dashboard;
