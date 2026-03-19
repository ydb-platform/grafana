SELECT
    u.id,
    u.uid,
    u.email,
    u.name,
    u.login,
    u.is_admin,
    u.is_disabled,
    u.last_seen_at,
    ua.auth_module,
    u.is_provisioned
FROM user AS u
LEFT JOIN (
    SELECT
        user_id,
        auth_module,
        ROW_NUMBER() OVER (PARTITION BY user_id ORDER BY created DESC) AS rn
    FROM user_auth
) ua ON u.id = ua.user_id
WHERE
    u.is_service_account = ?
    AND u.last_seen_at > ?
    AND ua.rn = 1
ORDER BY
    u.login ASC,
    u.email ASC
LIMIT 50;
