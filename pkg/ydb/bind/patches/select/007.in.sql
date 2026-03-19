SELECT `u`.`id`, `u`.`uid`, `u`.`email`, `u`.`name`, `u`.`login`, `u`.`is_admin`, `u`.`is_disabled`, `u`.`last_seen_at`, `user_auth`.`auth_module`, `u`.`is_provisioned` FROM `user` AS `u` LEFT JOIN user_auth ON user_auth.id=(
		SELECT id from user_auth
			WHERE user_auth.user_id = u.id
			ORDER BY user_auth.created DESC  LIMIT 1) WHERE (u.is_service_account = ? AND  1 = 1) ORDER BY `u`.`login` ASC, `u`.`email` ASC LIMIT 50
