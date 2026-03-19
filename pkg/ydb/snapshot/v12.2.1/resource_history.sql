CREATE TABLE IF NOT EXISTS `resource_history` (
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
    `generation` Int64,
    INDEX `IDX_resource_history_group_resource_resource_version` GLOBAL ASYNC ON (`group`, `resource`, `resource_version`),
    INDEX `IDX_resource_history_namespace_group_resource_action_version` GLOBAL ASYNC ON (`namespace`, `group`, `resource`, `action`, `resource_version`),
    INDEX `IDX_resource_history_namespace_group_resource_name_generation` GLOBAL ASYNC ON (`namespace`, `group`, `resource`, `name`, `generation`),
    INDEX `IDX_resource_history_resource_version` GLOBAL ASYNC ON (`resource_version`),
    INDEX `UQE_resource_history_namespace_group_name_version` GLOBAL UNIQUE SYNC ON (`namespace`, `group`, `resource`, `name`, `resource_version`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`guid`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);





