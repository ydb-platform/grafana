CREATE TABLE IF NOT EXISTS `dashboard_tag` (
  `id` Serial NOT NULL,
  `dashboard_id` Int64 NOT NULL,
  `term` Text NOT NULL,
  `dashboard_uid` Text,
  `org_id` Int64 DEFAULT 1,
  PRIMARY KEY (`id`),
  INDEX `IDX_dashboard_tag_dashboard_id` GLOBAL SYNC ON (`dashboard_id`),
  INDEX `IDX_dashboard_tag_dashboard_uid` GLOBAL SYNC ON (`dashboard_uid`)
);
