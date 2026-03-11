CREATE TABLE IF NOT EXISTS `alert_configuration_history` (
  `id` Serial NOT NULL,
  `org_id` Int64 NOT NULL DEFAULT 0,
  `alertmanager_configuration` Text NOT NULL,
  `configuration_hash` Text NOT NULL DEFAULT 'not-yet-calculated',
  `configuration_version` Text NOT NULL,
  `created_at` Int64 NOT NULL,
  `default` Int64 NOT NULL DEFAULT 0,
  `last_applied` Int64 NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
);
