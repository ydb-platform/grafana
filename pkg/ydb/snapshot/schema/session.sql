CREATE TABLE IF NOT EXISTS `session` (
    `key` Text,
    `data` String,
    `expiry` Uint32,
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`key`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);
