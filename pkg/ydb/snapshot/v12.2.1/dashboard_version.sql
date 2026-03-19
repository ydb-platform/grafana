CREATE TABLE IF NOT EXISTS `dashboard_version` (
    `id` Serial8 NOT NULL,
    `dashboard_id` Int64,
    `parent_version` Int64,
    `restored_from` Int64,
    `version` Int64,
    `created` Datetime64,
    `created_by` Int64,
    `message` Text,
    `data` Text,
    `api_version` Text,
    INDEX `IDX_dashboard_version_dashboard_id` GLOBAL ASYNC ON (`dashboard_id`),
    INDEX `UQE_dashboard_version_dashboard_id_version` GLOBAL UNIQUE SYNC ON (`dashboard_id`, `version`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);


