CREATE TABLE IF NOT EXISTS `org_user` (
    `id` Serial8 NOT NULL,
    `org_id` Int64,
    `user_id` Int64,
    `role` Text,
    `created` Datetime64,
    `updated` Datetime64,
    INDEX `IDX_org_user_org_id` GLOBAL ASYNC ON (`org_id`),
    INDEX `IDX_org_user_user_id` GLOBAL ASYNC ON (`user_id`),
    INDEX `UQE_org_user_org_id_user_id` GLOBAL UNIQUE SYNC ON (`org_id`, `user_id`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);



