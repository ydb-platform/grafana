CREATE TABLE IF NOT EXISTS `signing_key` (
    `id` Serial8 NOT NULL,
    `key_id` Text,
    `private_key` Text,
    `added_at` Datetime64,
    `expires_at` Datetime64,
    `alg` Text,
    INDEX `UQE_signing_key_key_id` GLOBAL UNIQUE SYNC ON (`key_id`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);

