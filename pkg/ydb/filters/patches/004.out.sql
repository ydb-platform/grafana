UPSERT INTO dashboard
SELECT
  folder.uid AS folder_uid,
  d.created AS created,
  d.data as data,
  d.org_id as org_id,
  d.slug as slug,
  d.title as title,
  d.updated as updated,
  d.version as version
FROM
  dashboard folder
  JOIN dashboard d ON d.folder_id = folder.id
WHERE
  d.is_folder = ?
