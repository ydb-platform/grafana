CREATE TABLE IF NOT EXISTS `user_auth` (
  `id` Serial NOT NULL,
  `user_id` Int64 NOT NULL,
  `auth_module` Text NOT NULL,
  `auth_id` Text NOT NULL,
  `created` Datetime64 NOT NULL,
  `o_auth_access_token` Text,
  `o_auth_refresh_token` Text,
  `o_auth_token_type` Text,
  `o_auth_expiry` Datetime64,
  `o_auth_id_token` Text,
  `external_uid` Text,
  PRIMARY KEY (`id`),
  INDEX `IDX_user_auth_auth_module_auth_id` GLOBAL SYNC ON (`auth_module`, `auth_id`),
  INDEX `IDX_user_auth_user_id` GLOBAL SYNC ON (`user_id`)
);
