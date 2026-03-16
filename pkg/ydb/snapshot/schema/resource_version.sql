CREATE TABLE IF NOT EXISTS `resource_version` (
    `group` Text,
    `resource` Text,
    `resource_version` Int64,
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`group`, `resource`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);
