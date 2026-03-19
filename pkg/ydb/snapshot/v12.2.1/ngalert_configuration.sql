CREATE TABLE IF NOT EXISTS `ngalert_configuration` (
    `id` Serial8 NOT NULL,
    `org_id` Int64,
    `alertmanagers` Text,
    `created_at` Int64,
    `updated_at` Int64,
    `send_alerts_to` Int64,
    INDEX `UQE_ngalert_configuration_org_id` GLOBAL UNIQUE SYNC ON (`org_id`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);

