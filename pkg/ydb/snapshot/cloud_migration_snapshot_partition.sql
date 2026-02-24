CREATE TABLE IF NOT EXISTS `cloud_migration_snapshot_partition` (
  `snapshot_uid` Text NOT NULL,
  `partition_number` Int64 NOT NULL,
  `resource_type` Text NOT NULL,
  `data` Bytes NOT NULL,
  PRIMARY KEY (
    `snapshot_uid`,
    `resource_type`,
    `partition_number`
  )
);
