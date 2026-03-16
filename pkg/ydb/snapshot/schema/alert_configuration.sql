CREATE TABLE IF NOT EXISTS `alert_configuration` (
    `id` Serial8 NOT NULL,
    `alertmanager_configuration` Text,
    `configuration_version` Text,
    `created_at` Int64,
    `default` Bool,
    `org_id` Int64,
    `configuration_hash` Text,
    INDEX `UQE_alert_configuration_org_id` GLOBAL UNIQUE SYNC ON (`org_id`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);

