DELETE FROM user_external_session WHERE NOT EXISTS (SELECT 1 FROM user_auth_token WHERE user_external_session.id = user_auth_token.external_session_id)
