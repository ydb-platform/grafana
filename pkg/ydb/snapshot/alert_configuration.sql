CREATE TABLE IF NOT EXISTS `alert_configuration` (
  `id` Serial NOT NULL,
  `alertmanager_configuration` Text NOT NULL,
  `configuration_version` Text NOT NULL,
  `created_at` Int64 NOT NULL,
  `default` Int64 NOT NULL DEFAULT 0,
  `org_id` Int64 NOT NULL DEFAULT 0,
  `configuration_hash` Text NOT NULL DEFAULT 'not-yet-calculated',
  PRIMARY KEY (`id`),
  INDEX `UQE_alert_configuration_org_id` GLOBAL UNIQUE SYNC ON (`org_id`)
);
