CREATE TABLE IF NOT EXISTS `cloud_migration_resource` (
    `id` Serial8 NOT NULL,
    `uid` Text,
    `resource_type` Text,
    `resource_uid` Text,
    `status` Text,
    `error_string` Text,
    `snapshot_uid` Text,
    `name` Text,
    `parent_name` Text,
    `error_code` Text,
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);
