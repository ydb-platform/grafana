CREATE TABLE IF NOT EXISTS `dashboard_provisioning` (
  `id` Serial NOT NULL,
  `dashboard_id` Int64,
  `name` Text NOT NULL,
  `external_id` Text NOT NULL,
  `updated` Int64 NOT NULL DEFAULT 0,
  `check_sum` Text,
  PRIMARY KEY (`id`),
  INDEX `IDX_dashboard_provisioning_dashboard_id` GLOBAL SYNC ON (`dashboard_id`),
  INDEX `IDX_dashboard_provisioning_dashboard_id_name` GLOBAL SYNC ON (`dashboard_id`, `name`)
);
