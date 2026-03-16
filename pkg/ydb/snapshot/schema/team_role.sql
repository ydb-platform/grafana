CREATE TABLE IF NOT EXISTS `team_role` (
    `id` Serial8 NOT NULL,
    `org_id` Int64,
    `team_id` Int64,
    `role_id` Int64,
    `created` Datetime64,
    INDEX `IDX_team_role_org_id` GLOBAL ASYNC ON (`org_id`),
    INDEX `IDX_team_role_team_id` GLOBAL ASYNC ON (`team_id`),
    INDEX `UQE_team_role_org_id_team_id_role_id` GLOBAL UNIQUE SYNC ON (`org_id`, `team_id`, `role_id`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);



