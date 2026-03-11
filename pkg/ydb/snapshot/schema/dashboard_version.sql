CREATE TABLE IF NOT EXISTS `dashboard_version` (
  `id` Serial NOT NULL,
  `dashboard_id` Int64 NOT NULL,
  `parent_version` Int64 NOT NULL,
  `restored_from` Int64 NOT NULL,
  `version` Int64 NOT NULL,
  `created` Datetime64 NOT NULL,
  `created_by` Int64 NOT NULL,
  `message` Text NOT NULL,
  `data` Text NOT NULL,
  `api_version` Text,
  PRIMARY KEY (`id`),
  INDEX `IDX_dashboard_version_dashboard_id` GLOBAL SYNC ON (`dashboard_id`),
  INDEX `UQE_dashboard_version_dashboard_id_version` GLOBAL UNIQUE SYNC ON (`dashboard_id`, `version`)
);
