CREATE TABLE IF NOT EXISTS `alert_rule_tag` (
    `id` Serial8 NOT NULL,
    `alert_id` Int64,
    `tag_id` Int64,
    INDEX `UQE_alert_rule_tag_alert_id_tag_id` GLOBAL UNIQUE SYNC ON (`alert_id`, `tag_id`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);

