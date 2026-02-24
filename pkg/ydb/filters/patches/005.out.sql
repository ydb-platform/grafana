UPSERT INTO dashboard
SELECT
  folder.parent_uid AS folder_uid,
  d.created AS created,
  d.data as data,
  d.org_id as org_id,
  d.slug as slug,
  d.title as title,
  d.updated as updated,
  d.version as version
FROM
  folder
  JOIN dashboard d ON d.org_id = folder.org_id
WHERE
  d.is_folder = ?
