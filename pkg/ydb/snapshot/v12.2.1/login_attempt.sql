CREATE TABLE IF NOT EXISTS `login_attempt` (
    `id` Serial8 NOT NULL,
    `username` Text,
    `ip_address` Text,
    `created` Int64,
    INDEX `IDX_login_attempt_username` GLOBAL ASYNC ON (`username`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);

