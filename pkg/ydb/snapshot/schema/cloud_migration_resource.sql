CREATE TABLE IF NOT EXISTS `cloud_migration_resource` (
  `id` Serial NOT NULL,
  `uid` Text NOT NULL,
  `resource_type` Text NOT NULL,
  `resource_uid` Text NOT NULL,
  `status` Text NOT NULL,
  `error_string` Text,
  `snapshot_uid` Text NOT NULL,
  `name` Text,
  `parent_name` Text,
  `error_code` Text,
  PRIMARY KEY (`id`)
);
