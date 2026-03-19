CREATE TABLE IF NOT EXISTS `alert_image` (
    `id` Serial8 NOT NULL,
    `token` Text,
    `path` Text,
    `url` Text,
    `created_at` Datetime64,
    `expires_at` Datetime64,
    INDEX `UQE_alert_image_token` GLOBAL UNIQUE SYNC ON (`token`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);

