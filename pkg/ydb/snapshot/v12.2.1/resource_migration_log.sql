CREATE TABLE IF NOT EXISTS `resource_migration_log` (
    `id` Serial8 NOT NULL,
    `migration_id` Text,
    `sql` Text,
    `success` Bool,
    `error` Text,
    `timestamp` Datetime64,
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);

