CREATE TABLE IF NOT EXISTS `alert_notification` (
  `id` Serial NOT NULL,
  `org_id` Int64 NOT NULL,
  `name` Text NOT NULL,
  `type` Text NOT NULL,
  `settings` Text NOT NULL,
  `created` Datetime64 NOT NULL,
  `updated` Datetime64 NOT NULL,
  `is_default` Uint8 NOT NULL DEFAULT 0,
  `frequency` Int64,
  `send_reminder` Int64 DEFAULT 0,
  `disable_resolve_message` Int64 NOT NULL DEFAULT 0,
  `uid` Text,
  `secure_settings` Text,
  PRIMARY KEY (`id`),
  INDEX `UQE_alert_notification_org_id_uid` GLOBAL UNIQUE SYNC ON (`org_id`, `uid`)
);
