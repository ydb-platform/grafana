UPDATE `user`
SET
    login = 'sa-' || Unwrap(CAST(org_id AS STRING)) || '-' || CASE
        WHEN SUBSTRING(login, 1, 3) == 'sa-' THEN SUBSTRING(login, 4)
        ELSE login
    END
WHERE login IS NOT NULL
AND is_service_account == 1
AND login NOT LIKE 'sa-' || Unwrap(CAST(org_id AS STRING)) || '-%';
