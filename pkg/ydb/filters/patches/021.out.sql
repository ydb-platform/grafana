UPSERT INTO `user`
SELECT
    u1.id AS id,
    u1.dest_login AS login
FROM (
    SELECT
        u.id AS id,
        'sa-' || CAST(u.org_id AS STRING) || SUBSTRING(
            u.login,
            CAST(LENGTH('sa-' || CAST(u.org_id AS STRING) || '-' || CAST(u.org_id AS STRING)) + 1 AS Uint32)
        ) AS dest_login
    FROM
        `user` AS u
    WHERE
        u.login IS NOT NULL
        AND u.is_service_account == 1
        AND u.org_id IS NOT NULL
        AND u.login LIKE Unwrap('sa-' || CAST(u.org_id AS STRING) || '-' || CAST(u.org_id AS STRING) || '-%')
) AS u1
LEFT JOIN
    `user` AS u2
ON
    u1.dest_login == u2.login
WHERE
    u2.login IS NULL
;
