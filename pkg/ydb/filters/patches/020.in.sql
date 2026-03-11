            UPDATE `user` AS u
            SET login = 'sa-' || CAST(u.org_id AS TEXT) || SUBSTRING(u.login, LENGTH('sa-'||CAST(u.org_id AS TEXT)||'-'||CAST(u.org_id AS TEXT))+1)
            WHERE u.login IS NOT NULL
                AND u.is_service_account = 1
                AND u.login LIKE 'sa-'||CAST(u.org_id AS TEXT)||'-'||CAST(u.org_id AS TEXT)||'-%'
                AND NOT EXISTS (
                    SELECT 1
                    FROM  `user`AS u2
                    WHERE u2.login = 'sa-' || CAST(u.org_id AS TEXT) || SUBSTRING(u.login, LENGTH('sa-'||CAST(u.org_id AS TEXT)||'-'||CAST(u.org_id AS TEXT))+1)
                );;
        
