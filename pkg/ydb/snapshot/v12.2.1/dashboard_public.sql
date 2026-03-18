CREATE TABLE IF NOT EXISTS `dashboard_public` (
    `uid` Text,
    `dashboard_uid` Text,
    `org_id` Int64,
    `time_settings` Text,
    `template_variables` Text,
    `access_token` Text,
    `created_by` Int64,
    `updated_by` Int64,
    `created_at` Datetime64,
    `updated_at` Datetime64,
    `is_enabled` Bool,
    `annotations_enabled` Bool,
    `time_selection_enabled` Bool,
    `share` Text,
    INDEX `IDX_dashboard_public_config_org_id_dashboard_uid` GLOBAL ASYNC ON (`org_id`, `dashboard_uid`),
    INDEX `UQE_dashboard_public_config_access_token` GLOBAL UNIQUE SYNC ON (`access_token`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`uid`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);


