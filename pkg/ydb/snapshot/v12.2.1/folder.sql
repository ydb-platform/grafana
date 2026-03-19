CREATE TABLE IF NOT EXISTS `folder` (
    `id` Serial8 NOT NULL,
    `uid` Text,
    `org_id` Int64,
    `title` Text,
    `description` Text,
    `parent_uid` Text,
    `created` Datetime64,
    `updated` Datetime64,
    INDEX `IDX_folder_org_id_parent_uid` GLOBAL ASYNC ON (`org_id`, `parent_uid`),
    INDEX `UQE_folder_org_id_uid` GLOBAL UNIQUE SYNC ON (`org_id`, `uid`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);


