CREATE TABLE IF NOT EXISTS `server_lock` (
    `id` Serial8 NOT NULL,
    `operation_uid` Text,
    `version` Int64,
    `last_execution` Int64,
    INDEX `UQE_server_lock_operation_uid` GLOBAL UNIQUE SYNC ON (`operation_uid`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);

