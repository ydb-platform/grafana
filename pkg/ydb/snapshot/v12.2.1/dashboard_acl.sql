CREATE TABLE IF NOT EXISTS `dashboard_acl` (
    `id` Serial8 NOT NULL,
    `org_id` Int64,
    `dashboard_id` Int64,
    `user_id` Int64,
    `team_id` Int64,
    `permission` Int64,
    `role` Text,
    `created` Datetime64,
    `updated` Datetime64,
    INDEX `IDX_dashboard_acl_org_id_role` GLOBAL ASYNC ON (`org_id`, `role`),
    INDEX `IDX_dashboard_acl_permission` GLOBAL ASYNC ON (`permission`),
    INDEX `IDX_dashboard_acl_team_id` GLOBAL ASYNC ON (`team_id`),
    INDEX `IDX_dashboard_acl_user_id` GLOBAL ASYNC ON (`user_id`),
    INDEX `UQE_dashboard_acl_dashboard_id_team_id` GLOBAL UNIQUE SYNC ON (`dashboard_id`, `team_id`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);





