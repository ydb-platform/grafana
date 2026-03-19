CREATE TABLE IF NOT EXISTS `sso_setting` (
    `id` Text,
    `provider` Text,
    `settings` Text,
    `created` Datetime64,
    `updated` Datetime64,
    `is_deleted` Bool,
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);
