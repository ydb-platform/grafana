CREATE TABLE IF NOT EXISTS `role` (
    `id` Serial8 NOT NULL,
    `name` Text,
    `description` Text,
    `version` Int64,
    `org_id` Int64,
    `uid` Text,
    `created` Datetime64,
    `updated` Datetime64,
    `display_name` Text,
    `group_name` Text,
    `hidden` Bool,
    INDEX `IDX_role_org_id` GLOBAL ASYNC ON (`org_id`),
    INDEX `UQE_role_org_id_name` GLOBAL UNIQUE SYNC ON (`org_id`, `name`),
    INDEX `UQE_role_uid` GLOBAL UNIQUE SYNC ON (`uid`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);



