CREATE TABLE IF NOT EXISTS `alert` (
    `id` Serial8 NOT NULL,
    `version` Int64,
    `dashboard_id` Int64,
    `panel_id` Int64,
    `org_id` Int64,
    `name` Text,
    `message` Text,
    `state` Text,
    `settings` Text,
    `frequency` Int64,
    `handler` Int64,
    `severity` Text,
    `silenced` Bool,
    `execution_error` Text,
    `eval_data` Text,
    `eval_date` Datetime64,
    `new_state_date` Datetime64,
    `state_changes` Int64,
    `created` Datetime64,
    `updated` Datetime64,
    `for` Int64,
    INDEX `IDX_alert_dashboard_id` GLOBAL ASYNC ON (`dashboard_id`),
    INDEX `IDX_alert_org_id_id` GLOBAL ASYNC ON (`org_id`, `id`),
    INDEX `IDX_alert_state` GLOBAL ASYNC ON (`state`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);



