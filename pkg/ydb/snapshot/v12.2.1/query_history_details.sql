CREATE TABLE IF NOT EXISTS `query_history_details` (
    `id` Serial8 NOT NULL,
    `query_history_item_uid` Text,
    `datasource_uid` Text,
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);
