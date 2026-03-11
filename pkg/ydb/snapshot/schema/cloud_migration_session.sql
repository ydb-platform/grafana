CREATE TABLE IF NOT EXISTS `cloud_migration_session` (
  `id` Serial NOT NULL,
  `uid` Text,
  `auth_token` Text,
  `slug` Text NOT NULL,
  `stack_id` Int64 NOT NULL,
  `region_slug` Text NOT NULL,
  `cluster_slug` Text NOT NULL,
  `created` Datetime64 NOT NULL,
  `updated` Datetime64 NOT NULL,
  `org_id` Int64 NOT NULL DEFAULT 1,
  PRIMARY KEY (`id`),
  INDEX `UQE_cloud_migration_session_uid` GLOBAL UNIQUE SYNC ON (`uid`)
);
