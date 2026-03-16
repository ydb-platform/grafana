CREATE TABLE IF NOT EXISTS `kv_store` (
    `id` Serial8 NOT NULL,
    `org_id` Int64,
    `namespace` Text,
    `key` Text,
    `value` Text,
    `created` Datetime64,
    `updated` Datetime64,
    INDEX `UQE_kv_store_org_id_namespace_key` GLOBAL UNIQUE SYNC ON (`org_id`, `namespace`, `key`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);

