CREATE TABLE IF NOT EXISTS `playlist` (
    `id` Serial8 NOT NULL,
    `name` Text,
    `interval` Text,
    `org_id` Int64,
    `uid` Text,
    `created_at` Int64,
    `updated_at` Int64,
    INDEX `UQE_playlist_org_id_uid` GLOBAL UNIQUE SYNC ON (`org_id`, `uid`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);

