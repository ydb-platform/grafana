CREATE TABLE IF NOT EXISTS `secrets` (
    `id` Serial8 NOT NULL,
    `org_id` Int64,
    `namespace` Text,
    `type` Text,
    `value` Text,
    `created` Datetime64,
    `updated` Datetime64,
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);

