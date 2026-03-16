CREATE TABLE IF NOT EXISTS `data_keys` (
    `name` Text,
    `active` Bool,
    `scope` Text,
    `provider` Text,
    `encrypted_data` String,
    `created` Datetime64,
    `updated` Datetime64,
    `label` Text,
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`name`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);
