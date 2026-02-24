UPDATE library_element ON
SELECT
  library_element.id AS id,
  dashboard.uid AS folder_uid
FROM
  library_element
  JOIN dashboard ON dashboard.id = library_element.folder_id
  AND library_element.org_id = dashboard.org_id
