CREATE TABLE IF NOT EXISTS `dashboard` (
    `id` Serial8 NOT NULL,
    `version` Int64,
    `slug` Text,
    `title` Text,
    `data` Text,
    `org_id` Int64,
    `created` Datetime64,
    `updated` Datetime64,
    `updated_by` Int64,
    `created_by` Int64,
    `gnet_id` Int64,
    `plugin_id` Text,
    `folder_id` Int64,
    `is_folder` Bool,
    `has_acl` Bool,
    `uid` Text,
    `is_public` Bool,
    `folder_uid` Text,
    `deleted` Datetime64,
    `api_version` Text,
    INDEX `IDX_dashboard_deleted` GLOBAL ASYNC ON (`deleted`),
    INDEX `IDX_dashboard_gnet_id` GLOBAL ASYNC ON (`gnet_id`),
    INDEX `IDX_dashboard_is_folder` GLOBAL ASYNC ON (`is_folder`),
    INDEX `IDX_dashboard_org_id` GLOBAL ASYNC ON (`org_id`),
    INDEX `IDX_dashboard_org_id_folder_id_title` GLOBAL ASYNC ON (`org_id`, `folder_id`, `title`),
    INDEX `IDX_dashboard_org_id_is_folder_deleted` GLOBAL SYNC ON (`org_id`, `is_folder`, `deleted`),
    INDEX `IDX_dashboard_org_id_is_folder_id` GLOBAL SYNC ON (`org_id`, `is_folder`, `id`),
    INDEX `IDX_dashboard_org_id_plugin_id` GLOBAL ASYNC ON (`org_id`, `plugin_id`),
    INDEX `IDX_dashboard_title` GLOBAL ASYNC ON (`title`),
    INDEX `UQE_dashboard_org_id_uid` GLOBAL UNIQUE SYNC ON (`org_id`, `uid`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);










