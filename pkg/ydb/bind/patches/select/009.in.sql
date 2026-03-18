SELECT COALESCE(ou.role, 'None') AS role, u.is_admin
FROM `user` as u
  LEFT JOIN `org_user` as ou ON ou.user_id = u.id AND ou.org_id = ?
WHERE u.id = ?
