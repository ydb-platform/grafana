CREATE TABLE IF NOT EXISTS `permission` (
    `id` Serial8 NOT NULL,
    `role_id` Int64,
    `action` Text,
    `scope` Text,
    `created` Datetime64,
    `updated` Datetime64,
    `kind` Text,
    `attribute` Text,
    `identifier` Text,
    INDEX `IDX_permission_identifier` GLOBAL ASYNC ON (`identifier`),
    INDEX `IDX_permission_role_id_action` GLOBAL ASYNC ON (`role_id`, `action`),
    INDEX `UQE_permission_action_scope_role_id` GLOBAL UNIQUE SYNC ON (`action`, `scope`, `role_id`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);



