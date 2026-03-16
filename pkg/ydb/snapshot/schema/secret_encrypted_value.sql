CREATE TABLE IF NOT EXISTS `secret_encrypted_value` (
    `namespace` Text,
    `name` Text,
    `version` Int64,
    `encrypted_data` String,
    `created` Int64,
    `updated` Int64,
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`namespace`, `name`, `version`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);
