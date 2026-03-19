CREATE TABLE IF NOT EXISTS `team_member` (
    `id` Serial8 NOT NULL,
    `org_id` Int64,
    `team_id` Int64,
    `user_id` Int64,
    `created` Datetime64,
    `updated` Datetime64,
    `external` Bool,
    `permission` Int64,
    INDEX `IDX_team_member_org_id` GLOBAL ASYNC ON (`org_id`),
    INDEX `IDX_team_member_team_id` GLOBAL ASYNC ON (`team_id`),
    INDEX `IDX_team_member_user_id_org_id` GLOBAL ASYNC ON (`user_id`, `org_id`),
    INDEX `UQE_team_member_org_id_team_id_user_id` GLOBAL UNIQUE SYNC ON (`org_id`, `team_id`, `user_id`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);




