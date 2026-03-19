CREATE TABLE IF NOT EXISTS `alert_rule_version` (
    `id` Serial8 NOT NULL,
    `rule_org_id` Int64,
    `rule_uid` Text,
    `rule_namespace_uid` Text,
    `rule_group` Text,
    `parent_version` Int64,
    `restored_from` Int64,
    `version` Int64,
    `created` Datetime64,
    `title` Text,
    `condition` Text,
    `data` Text,
    `interval_seconds` Int64,
    `no_data_state` Text,
    `exec_err_state` Text,
    `for` Int64,
    `annotations` Text,
    `labels` Text,
    `rule_group_idx` Int64,
    `is_paused` Bool,
    `notification_settings` Text,
    `record` Text,
    `metadata` Text,
    `created_by` Text,
    `rule_guid` Text,
    `keep_firing_for` Int64,
    `missing_series_evals_to_resolve` Int64,
    INDEX `IDX_alert_rule_version_rule_org_id_rule_namespace_uid_rule_group` GLOBAL ASYNC ON (`rule_org_id`, `rule_namespace_uid`, `rule_group`),
    INDEX `UQE_alert_rule_version_rule_guid_version` GLOBAL UNIQUE SYNC ON (`rule_guid`, `version`),
    INDEX `UQE_alert_rule_version_rule_org_id_rule_uid_rule_guid_version` GLOBAL UNIQUE SYNC ON (`rule_org_id`, `rule_uid`, `rule_guid`, `version`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);



