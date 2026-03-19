CREATE TABLE IF NOT EXISTS `cloud_migration_snapshot` (
    `id` Serial8 NOT NULL,
    `uid` Text,
    `session_uid` Text,
    `created` Datetime64,
    `updated` Datetime64,
    `finished` Datetime64,
    `upload_url` Text,
    `status` Text,
    `local_directory` Text,
    `gms_snapshot_uid` Text,
    `encryption_key` Text,
    `error_string` Text,
    `resource_storage_type` Text,
    `encryption_algo` Text,
    `metadata` String,
    `public_key` String,
    INDEX `UQE_cloud_migration_snapshot_uid` GLOBAL UNIQUE SYNC ON (`uid`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);

