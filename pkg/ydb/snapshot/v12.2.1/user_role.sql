CREATE TABLE IF NOT EXISTS `user_role` (
    `id` Serial8 NOT NULL,
    `org_id` Int64,
    `user_id` Int64,
    `role_id` Int64,
    `created` Datetime64,
    `group_mapping_uid` Text,
    INDEX `IDX_user_role_org_id` GLOBAL ASYNC ON (`org_id`),
    INDEX `IDX_user_role_user_id` GLOBAL ASYNC ON (`user_id`),
    INDEX `UQE_user_role_org_id_user_id_role_id_group_mapping_uid` GLOBAL UNIQUE SYNC ON (`org_id`, `user_id`, `role_id`, `group_mapping_uid`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);



