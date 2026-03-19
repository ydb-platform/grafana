CREATE TABLE IF NOT EXISTS `query_history_star` (
    `id` Serial8 NOT NULL,
    `query_uid` Text,
    `user_id` Int64,
    `org_id` Int64,
    INDEX `UQE_query_history_star_user_id_query_uid` GLOBAL UNIQUE SYNC ON (`user_id`, `query_uid`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);

