CREATE TABLE IF NOT EXISTS `cloud_migration_snapshot_partition` (
    `snapshot_uid` Text,
    `partition_number` Int64,
    `resource_type` Text,
    `data` Text,
    INDEX `UQE_cloud_migration_snapshot_partition_srp_unique` GLOBAL UNIQUE SYNC ON (`snapshot_uid`, `resource_type`, `partition_number`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`snapshot_uid`, `partition_number`, `resource_type`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);

