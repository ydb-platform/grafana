CREATE TABLE IF NOT EXISTS `builtin_role` (
    `id` Serial8 NOT NULL,
    `role` Text,
    `role_id` Int64,
    `created` Datetime64,
    `updated` Datetime64,
    `org_id` Int64,
    INDEX `IDX_builtin_role_org_id` GLOBAL ASYNC ON (`org_id`),
    INDEX `IDX_builtin_role_role` GLOBAL ASYNC ON (`role`),
    INDEX `IDX_builtin_role_role_id` GLOBAL ASYNC ON (`role_id`),
    INDEX `UQE_builtin_role_org_id_role_id_role` GLOBAL UNIQUE SYNC ON (`org_id`, `role_id`, `role`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);




