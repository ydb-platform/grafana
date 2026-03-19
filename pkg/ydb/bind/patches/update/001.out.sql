$to_update =
  SELECT guid FROM (
    SELECT
      guid,
      ROW_NUMBER() OVER (ORDER BY created ASC) AS rn
    FROM `secret_secure_value`
    WHERE
      `active` = FALSE AND
      ? - `created` > ? AND
      ? - `lease_created` > ?
  ) AS sub
  WHERE rn <= ?
;
UPDATE
  `secret_secure_value`
SET
  `lease_token` = ?,
  `lease_created` = ?
WHERE guid IN (SELECT guid FROM $to_update)
;
