CREATE TABLE IF NOT EXISTS `secret_data_key` (
    `uid` Text,
    `namespace` Text,
    `label` Text,
    `active` Bool,
    `provider` Text,
    `encrypted_data` String,
    `created` Datetime64,
    `updated` Datetime64,
    INDEX `IDX_secret_data_key_namespace_label_active` GLOBAL ASYNC ON (`namespace`, `label`, `active`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`uid`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);

