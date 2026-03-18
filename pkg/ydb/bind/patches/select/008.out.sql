SELECT
  team.id AS id,
  team.uid AS uid,
  team.org_id AS org_id,
  team.name AS name,
  team.email AS email,
  team.external_uid AS external_uid,
  team.is_provisioned AS is_provisioned,
  COUNT(member.team_id) AS member_count
FROM
  team AS team
    LEFT JOIN team_member AS member ON member.team_id == team.id
WHERE
  team.org_id == ?
GROUP BY
  team.id,
  team.uid,
  team.org_id,
  team.name,
  team.email,
  team.external_uid,
  team.is_provisioned
ORDER BY
  name ASC
  LIMIT
  30 OFFSET 0;
