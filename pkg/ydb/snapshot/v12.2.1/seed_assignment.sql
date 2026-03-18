CREATE TABLE IF NOT EXISTS `seed_assignment` (
    `id` Serial8 NOT NULL,
    `builtin_role` Text,
    `role_name` Text,
    `action` Text,
    `scope` Text,
    `origin` Text,
    INDEX `UQE_seed_assignment_builtin_role_action_scope` GLOBAL UNIQUE SYNC ON (`builtin_role`, `action`, `scope`),
    INDEX `UQE_seed_assignment_builtin_role_role_name` GLOBAL UNIQUE SYNC ON (`builtin_role`, `role_name`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);


