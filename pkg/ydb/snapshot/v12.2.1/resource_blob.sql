CREATE TABLE IF NOT EXISTS `resource_blob` (
    `uuid` Text,
    `created` Datetime64,
    `group` Text,
    `resource` Text,
    `namespace` Text,
    `name` Text,
    `value` Text,
    `hash` Text,
    `content_type` Text,
    INDEX `IDX_resource_blob_created` GLOBAL ASYNC ON (`created`),
    INDEX `IDX_resource_history_namespace_group_name` GLOBAL ASYNC ON (`namespace`, `group`, `resource`, `name`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`uuid`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);


