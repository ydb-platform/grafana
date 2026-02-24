SELECT
    SUM(serviceaccounts) AS serviceaccounts,
    SUM(serviceaccount_tokens) AS serviceaccount_tokens,
    SUM(serviceaccounts_with_no_role) AS serviceaccounts_with_no_role
FROM (
    SELECT *
    FROM
        (
            SELECT COUNT(*) AS serviceaccounts
            FROM `user`
            WHERE is_service_account = true
        ) UNION ALL (
            SELECT COUNT(*) AS serviceaccount_tokens
            FROM `api_key`
            WHERE service_account_id IS NOT NULL
        ) UNION ALL (
            SELECT COUNT(*) AS serviceaccounts_with_no_role
            FROM `org_user` AS ou
            JOIN `user` AS u
            ON u.id = ou.user_id
            WHERE
                u.is_disabled = false
                AND u.is_service_account = true
                AND ou.role='admin'
        )
)
