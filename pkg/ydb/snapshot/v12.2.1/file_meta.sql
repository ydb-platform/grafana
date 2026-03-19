CREATE TABLE IF NOT EXISTS `file_meta` (
    `path_hash` Text,
    `key` Text,
    `value` Text,
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`path_hash`, `key`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);
