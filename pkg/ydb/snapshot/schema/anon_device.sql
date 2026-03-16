CREATE TABLE IF NOT EXISTS `anon_device` (
    `id` Serial8 NOT NULL,
    `client_ip` Text,
    `created_at` Datetime64,
    `device_id` Text,
    `updated_at` Datetime64,
    `user_agent` Text,
    INDEX `IDX_anon_device_updated_at` GLOBAL ASYNC ON (`updated_at`),
    INDEX `UQE_anon_device_device_id` GLOBAL UNIQUE SYNC ON (`device_id`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);


