SELECT
  dashboard.id,
  dashboard.org_id,
  dashboard.uid,
  dashboard.title,
  dashboard.slug,
  dashboard_tag.term,
  dashboard.is_folder,
  dashboard.folder_id,
  dashboard.deleted,
  folder.uid AS folder_uid,
  folder.title AS folder_slug,
  folder.title AS folder_title,
  dashboard.title AS title
FROM (
  SELECT dashboard.id AS id, dashboard.title AS title
  FROM dashboard
  WHERE dashboard.org_id = $1
    AND dashboard.uid IN $2
    AND dashboard.is_folder = true
    AND dashboard.folder_uid IS NULL
    AND dashboard.uid != $3
    AND (dashboard.folder_uid != $4 OR dashboard.folder_uid IS NULL)
  ORDER BY title ASC
  LIMIT 50 OFFSET 0
) AS ids
INNER JOIN dashboard ON ids.id = dashboard.id
LEFT OUTER JOIN folder ON folder.uid = dashboard.folder_uid AND folder.org_id = dashboard.org_id
LEFT OUTER JOIN dashboard_tag ON dashboard.id = dashboard_tag.dashboard_id
ORDER BY dashboard.title ASC
