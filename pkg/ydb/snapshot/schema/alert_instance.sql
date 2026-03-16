CREATE TABLE IF NOT EXISTS `alert_instance` (
    `rule_org_id` Int64,
    `rule_uid` Text,
    `labels` Text,
    `labels_hash` Text,
    `current_state` Text,
    `current_state_since` Int64,
    `last_eval_time` Int64,
    `current_state_end` Int64,
    `current_reason` Text,
    `result_fingerprint` Text,
    `resolved_at` Int64,
    `last_sent_at` Int64,
    `fired_at` Int64,
    INDEX `IDX_alert_instance_rule_org_id_current_state` GLOBAL ASYNC ON (`rule_org_id`, `current_state`),
    INDEX `IDX_alert_instance_rule_org_id_rule_uid_current_state` GLOBAL ASYNC ON (`rule_org_id`, `rule_uid`, `current_state`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`rule_org_id`, `rule_uid`, `labels_hash`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);


