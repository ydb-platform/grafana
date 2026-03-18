CREATE TABLE IF NOT EXISTS `secret_keeper` (
    `guid` Text,
    `name` Text,
    `namespace` Text,
    `annotations` Text,
    `labels` Text,
    `created` Int64,
    `created_by` Text,
    `updated` Int64,
    `updated_by` Text,
    `description` Text,
    `type` Text,
    `payload` Text,
    INDEX `UQE_secret_keeper_namespace_name` GLOBAL UNIQUE SYNC ON (`namespace`, `name`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`guid`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);

