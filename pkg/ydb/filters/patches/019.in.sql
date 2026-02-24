            UPDATE `user`
            SET login = 'sa-' || CAST(org_id AS TEXT) || '-' ||
                CASE
                    WHEN SUBSTR(login, 1, 3) = 'sa-' THEN SUBSTR(login, 4)
                    ELSE login
                END
            WHERE login IS NOT NULL
              AND is_service_account = 1
              AND login NOT LIKE 'sa-' || CAST(org_id AS TEXT) || '-%';
        
