CREATE TABLE IF NOT EXISTS `user_external_session` (
  `id` Serial NOT NULL,
  `user_auth_id` Int64 NOT NULL,
  `user_id` Int64 NOT NULL,
  `auth_module` Text NOT NULL,
  `access_token` Text,
  `id_token` Text,
  `refresh_token` Text,
  `session_id` Text,
  `session_id_hash` Text,
  `name_id` Text,
  `name_id_hash` Text,
  `expires_at` Datetime64,
  `created_at` Datetime64 NOT NULL,
  PRIMARY KEY (`id`)
);
