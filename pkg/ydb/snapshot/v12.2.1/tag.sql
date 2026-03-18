CREATE TABLE IF NOT EXISTS `tag` (
    `id` Serial8 NOT NULL,
    `key` Text,
    `value` Text,
    INDEX `UQE_tag_key_value` GLOBAL UNIQUE SYNC ON (`key`, `value`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);

