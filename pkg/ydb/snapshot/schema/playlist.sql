CREATE TABLE IF NOT EXISTS `playlist` (
  `id` Serial NOT NULL,
  `name` Text NOT NULL,
  `interval` Text NOT NULL,
  `org_id` Int64 NOT NULL,
  `created_at` Int64 NOT NULL DEFAULT 0,
  `updated_at` Int64 NOT NULL DEFAULT 0,
  `uid` Text NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `UQE_playlist_org_id_uid` GLOBAL UNIQUE SYNC ON (`org_id`, `uid`)
);
