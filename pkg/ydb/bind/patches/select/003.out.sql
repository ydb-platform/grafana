SELECT
    role AS bitrole,
    active,
    COUNT(role) AS count
FROM (
    SELECT
        last_seen_at > ? AS active,
        last_seen_at > ? AS daily_active,
        SUM(role) AS role
    FROM (
        SELECT
            u.id AS id,
            CASE org_user.role
                WHEN 'Admin' THEN 4
                WHEN 'Editor' THEN 2
                ELSE 1
            END AS role,
            u.last_seen_at AS last_seen_at
        FROM
            `user` AS u
        INNER JOIN
            org_user
        ON
            org_user.user_id == u.id
        GROUP BY
            u.id,
            u.last_seen_at,
            org_user.role
    ) AS t2
    GROUP BY
        id,
        last_seen_at
) AS t1
GROUP BY
    active,
    daily_active,
    role
;
