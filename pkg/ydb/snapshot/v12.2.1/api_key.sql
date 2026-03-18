CREATE TABLE IF NOT EXISTS `api_key` (
    `id` Serial8 NOT NULL,
    `org_id` Int64,
    `name` Text,
    `key` Text,
    `role` Text,
    `created` Datetime64,
    `updated` Datetime64,
    `expires` Int64,
    `service_account_id` Int64,
    `last_used_at` Datetime64,
    `is_revoked` Bool,
    INDEX `IDX_api_key_org_id` GLOBAL ASYNC ON (`org_id`),
    INDEX `UQE_api_key_key` GLOBAL UNIQUE SYNC ON (`key`),
    INDEX `UQE_api_key_org_id_name` GLOBAL UNIQUE SYNC ON (`org_id`, `name`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);



