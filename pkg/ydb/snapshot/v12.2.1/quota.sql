CREATE TABLE IF NOT EXISTS `quota` (
    `id` Serial8 NOT NULL,
    `org_id` Int64,
    `user_id` Int64,
    `target` Text,
    `limit` Int64,
    `created` Datetime64,
    `updated` Datetime64,
    INDEX `UQE_quota_org_id_user_id_target` GLOBAL UNIQUE SYNC ON (`org_id`, `user_id`, `target`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);

