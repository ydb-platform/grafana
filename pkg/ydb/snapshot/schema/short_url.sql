CREATE TABLE IF NOT EXISTS `short_url` (
    `id` Serial8 NOT NULL,
    `org_id` Int64,
    `uid` Text,
    `path` Text,
    `created_by` Int64,
    `created_at` Int64,
    `last_seen_at` Int64,
    INDEX `UQE_short_url_org_id_uid` GLOBAL UNIQUE SYNC ON (`org_id`, `uid`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);

