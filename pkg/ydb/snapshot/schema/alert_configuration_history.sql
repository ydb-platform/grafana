CREATE TABLE IF NOT EXISTS `alert_configuration_history` (
    `id` Serial8 NOT NULL,
    `org_id` Int64,
    `alertmanager_configuration` Text,
    `configuration_hash` Text,
    `configuration_version` Text,
    `created_at` Int64,
    `default` Bool,
    `last_applied` Int64,
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);

