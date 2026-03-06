CREATE TABLE IF NOT EXISTS `query_history` (
  `id` Serial NOT NULL,
  `uid` Text NOT NULL,
  `org_id` Int64 NOT NULL,
  `datasource_uid` Text NOT NULL,
  `created_by` Int64 NOT NULL,
  `created_at` Int64 NOT NULL,
  `comment` Text NOT NULL,
  `queries` Text NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `IDX_query_history_org_id_created_by_datasource_uid` GLOBAL SYNC ON (`org_id`, `created_by`, `datasource_uid`)
);
