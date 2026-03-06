CREATE TABLE IF NOT EXISTS `sso_setting` (
  `id` Text NOT NULL,
  PRIMARY KEY (`id`),
  `provider` Text NOT NULL,
  `settings` Text NOT NULL,
  `created` Datetime64 NOT NULL,
  `updated` Datetime64 NOT NULL,
  `is_deleted` Bool NOT NULL DEFAULT false
);
