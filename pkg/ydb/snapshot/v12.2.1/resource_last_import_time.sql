CREATE TABLE IF NOT EXISTS `resource_last_import_time` (
    `group` Text,
    `resource` Text,
    `namespace` Text,
    `last_import_time` Datetime64,
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`group`, `resource`, `namespace`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = DISABLED
);
