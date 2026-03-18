CREATE TABLE IF NOT EXISTS `resource` (
    `guid` Text,
    `resource_version` Int64,
    `group` Text,
    `resource` Text,
    `namespace` Text,
    `name` Text,
    `value` Text,
    `action` Int64,
    `label_set` Text,
    `previous_resource_version` Int64,
    `folder` Text,
    INDEX `IDX_resource_group_resource` GLOBAL ASYNC ON (`group`, `resource`),
    INDEX `UQE_resource_namespace_group_resource_name` GLOBAL UNIQUE SYNC ON (`namespace`, `group`, `resource`, `name`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`guid`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);


