CREATE TABLE IF NOT EXISTS `alert_notification` (
    `id` Serial8 NOT NULL,
    `org_id` Int64,
    `name` Text,
    `type` Text,
    `settings` Text,
    `created` Datetime64,
    `updated` Datetime64,
    `is_default` Bool,
    `frequency` Int64,
    `send_reminder` Bool,
    `disable_resolve_message` Bool,
    `uid` Text,
    `secure_settings` Text,
    INDEX `UQE_alert_notification_org_id_uid` GLOBAL UNIQUE SYNC ON (`org_id`, `uid`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);

