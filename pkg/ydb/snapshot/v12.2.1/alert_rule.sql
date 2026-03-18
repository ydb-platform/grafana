CREATE TABLE IF NOT EXISTS `alert_rule` (
    `id` Serial8 NOT NULL,
    `org_id` Int64,
    `title` Text,
    `condition` Text,
    `data` Text,
    `updated` Datetime64,
    `interval_seconds` Int64,
    `version` Int64,
    `uid` Text,
    `namespace_uid` Text,
    `rule_group` Text,
    `no_data_state` Text,
    `exec_err_state` Text,
    `for` Int64,
    `annotations` Text,
    `labels` Text,
    `dashboard_uid` Text,
    `panel_id` Int64,
    `rule_group_idx` Int64,
    `is_paused` Bool,
    `notification_settings` Text,
    `record` Text,
    `metadata` Text,
    `updated_by` Text,
    `guid` Text,
    `keep_firing_for` Int64,
    `missing_series_evals_to_resolve` Int64,
    INDEX `IDX_alert_rule_org_id_dashboard_uid_panel_id` GLOBAL ASYNC ON (`org_id`, `dashboard_uid`, `panel_id`),
    INDEX `IDX_alert_rule_org_id_namespace_uid_rule_group` GLOBAL ASYNC ON (`org_id`, `namespace_uid`, `rule_group`),
    INDEX `UQE_alert_rule_guid` GLOBAL UNIQUE SYNC ON (`guid`),
    INDEX `UQE_alert_rule_org_id_uid` GLOBAL UNIQUE SYNC ON (`org_id`, `uid`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);




