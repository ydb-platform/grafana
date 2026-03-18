CREATE TABLE IF NOT EXISTS `dashboard_snapshot` (
    `id` Serial8 NOT NULL,
    `name` Text,
    `key` Text,
    `delete_key` Text,
    `org_id` Int64,
    `user_id` Int64,
    `external` Bool,
    `external_url` Text,
    `dashboard` Text,
    `expires` Datetime64,
    `created` Datetime64,
    `updated` Datetime64,
    `external_delete_url` Text,
    `dashboard_encrypted` Text,
    INDEX `IDX_dashboard_snapshot_user_id` GLOBAL ASYNC ON (`user_id`),
    INDEX `UQE_dashboard_snapshot_delete_key` GLOBAL UNIQUE SYNC ON (`delete_key`),
    INDEX `UQE_dashboard_snapshot_key` GLOBAL UNIQUE SYNC ON (`key`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);



