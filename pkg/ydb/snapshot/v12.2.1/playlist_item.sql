CREATE TABLE IF NOT EXISTS `playlist_item` (
    `id` Serial8 NOT NULL,
    `playlist_id` Int64,
    `type` Text,
    `value` Text,
    `title` Text,
    `order` Int64,
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);

