CREATE TABLE IF NOT EXISTS `ngalert_configuration` (
  `id` Serial NOT NULL,
  `org_id` Int64 NOT NULL,
  `alertmanagers` Text,
  `created_at` Int64 NOT NULL,
  `updated_at` Int64 NOT NULL,
  `send_alerts_to` Int64 NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  INDEX `UQE_ngalert_configuration_org_id` GLOBAL UNIQUE SYNC ON (`org_id`)
);
