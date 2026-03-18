SELECT
		u.id                  as user_id,
		u.uid                 as user_uid,
		u.is_admin            as is_grafana_admin,
		u.email               as email,
		u.email_verified      as email_verified,
		u.login               as login,
		u.name                as name,
		u.is_disabled         as is_disabled,
		u.help_flags1         as help_flags1,
		u.last_seen_at        as last_seen_at,
		org.name              as org_name,
		org_user.role         as org_role,
		org.id                as org_id,
		u.is_service_account  as is_service_account
		FROM `user` as u
		LEFT OUTER JOIN org_user on org_user.org_id = 1 and org_user.user_id = u.id
		LEFT OUTER JOIN org on org.id = org_user.org_id WHERE u.id=?
