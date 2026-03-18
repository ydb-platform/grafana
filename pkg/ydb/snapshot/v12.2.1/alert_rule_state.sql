CREATE TABLE IF NOT EXISTS `alert_rule_state` (
    `id` Serial8 NOT NULL,
    `org_id` Int64,
    `rule_uid` Text,
    `data` Text,
    `updated_at` Datetime64,
    INDEX `UQE_alert_rule_state_org_id_rule_uid` GLOBAL UNIQUE SYNC ON (`org_id`, `rule_uid`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);

