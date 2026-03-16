DELETE FROM folder ON
SELECT
  folder.id AS id
FROM
  folder
  JOIN dashboard ON dashboard.uid = folder.uid
WHERE
  dashboard.org_id = folder.org_id
  AND dashboard.is_folder = true
