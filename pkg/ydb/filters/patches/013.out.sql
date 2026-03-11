UPDATE dashboard_tag ON
SELECT
  dashboard.uid AS dashboard_uid,
  dashboard.org_id AS org_id,
  dashboard_tag.id AS id
FROM
  dashboard_tag
  LEFT JOIN dashboard ON dashboard_tag.dashboard_id = dashboard.id
WHERE
  dashboard_tag.dashboard_uid IS NULL
  OR dashboard_tag.org_id IS NULL;
