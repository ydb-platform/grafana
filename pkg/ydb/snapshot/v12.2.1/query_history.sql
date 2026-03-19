CREATE TABLE IF NOT EXISTS `query_history` (
    `id` Serial8 NOT NULL,
    `uid` Text,
    `org_id` Int64,
    `datasource_uid` Text,
    `created_by` Int64,
    `created_at` Int64,
    `comment` Text,
    `queries` Text,
    INDEX `IDX_query_history_org_id_created_by_datasource_uid` GLOBAL ASYNC ON (`org_id`, `created_by`, `datasource_uid`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);

