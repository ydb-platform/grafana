DELETE FROM user_external_session
WHERE
  id IN (
    SELECT
      ues.id
    FROM
      user_external_session ues
      LEFT JOIN user_auth_token uat ON ues.id = uat.external_session_id
    WHERE
      uat.external_session_id IS NULL
  )
