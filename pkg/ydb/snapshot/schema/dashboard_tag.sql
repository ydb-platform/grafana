CREATE TABLE IF NOT EXISTS `dashboard_tag` (
    `id` Serial8 NOT NULL,
    `dashboard_id` Int64,
    `term` Text,
    `dashboard_uid` Text,
    `org_id` Int64,
    INDEX `IDX_dashboard_tag_dashboard_id` GLOBAL ASYNC ON (`dashboard_id`),
    INDEX `IDX_dashboard_tag_dashboard_uid` GLOBAL ASYNC ON (`dashboard_uid`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);


