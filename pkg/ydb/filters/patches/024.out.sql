UPSERT INTO
    alert_rule_version
SELECT
  arv.id AS id,
  ar.guid AS rule_guid
FROM
  alert_rule_version AS arv
    JOIN
  alert_rule AS ar
  ON
    ar.uid == arv.rule_uid
    AND ar.org_id == arv.rule_org_id
;
