CREATE TABLE IF NOT EXISTS `alert_notification_state` (
    `id` Serial8 NOT NULL,
    `org_id` Int64,
    `alert_id` Int64,
    `notifier_id` Int64,
    `state` Text,
    `version` Int64,
    `updated_at` Int64,
    `alert_rule_state_updated_version` Int64,
    INDEX `IDX_alert_notification_state_alert_id` GLOBAL ASYNC ON (`alert_id`),
    INDEX `UQE_alert_notification_state_org_id_alert_id_notifier_id` GLOBAL UNIQUE SYNC ON (`org_id`, `alert_id`, `notifier_id`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);


