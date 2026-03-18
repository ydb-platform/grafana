CREATE TABLE IF NOT EXISTS `annotation` (
    `id` Serial8 NOT NULL,
    `org_id` Int64,
    `alert_id` Int64,
    `user_id` Int64,
    `dashboard_id` Int64,
    `panel_id` Int64,
    `category_id` Int64,
    `type` Text,
    `title` Text,
    `text` Text,
    `metric` Text,
    `prev_state` Text,
    `new_state` Text,
    `data` Text,
    `epoch` Int64,
    `region_id` Int64,
    `tags` Text,
    `created` Int64,
    `updated` Int64,
    `epoch_end` Int64,
    `dashboard_uid` Text,
    INDEX `IDX_annotation_alert_id` GLOBAL ASYNC ON (`alert_id`),
    INDEX `IDX_annotation_org_id_alert_id` GLOBAL ASYNC ON (`org_id`, `alert_id`),
    INDEX `IDX_annotation_org_id_created` GLOBAL ASYNC ON (`org_id`, `created`),
    INDEX `IDX_annotation_org_id_dashboard_id_epoch_end_epoch` GLOBAL ASYNC ON (`org_id`, `dashboard_id`, `epoch_end`, `epoch`),
    INDEX `IDX_annotation_org_id_epoch_end_epoch` GLOBAL ASYNC ON (`org_id`, `epoch_end`, `epoch`),
    INDEX `IDX_annotation_org_id_type` GLOBAL ASYNC ON (`org_id`, `type`),
    INDEX `IDX_annotation_org_id_updated` GLOBAL ASYNC ON (`org_id`, `updated`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);







