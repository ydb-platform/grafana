CREATE TABLE IF NOT EXISTS `dashboard_provisioning` (
    `id` Serial8 NOT NULL,
    `dashboard_id` Int64,
    `name` Text,
    `external_id` Text,
    `updated` Int64,
    `check_sum` Text,
    INDEX `IDX_dashboard_provisioning_dashboard_id` GLOBAL ASYNC ON (`dashboard_id`),
    INDEX `IDX_dashboard_provisioning_dashboard_id_name` GLOBAL ASYNC ON (`dashboard_id`, `name`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);


