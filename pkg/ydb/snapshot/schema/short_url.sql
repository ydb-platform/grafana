CREATE TABLE IF NOT EXISTS `short_url` (
  `id` Serial NOT NULL,
  `org_id` Int64 NOT NULL,
  `uid` Text NOT NULL,
  `path` Text NOT NULL,
  `created_by` Int64 NOT NULL,
  `created_at` Int64 NOT NULL,
  `last_seen_at` Int64,
  PRIMARY KEY (`id`),
  INDEX `UQE_short_url_org_id_uid` GLOBAL UNIQUE SYNC ON (`org_id`, `uid`)
);
