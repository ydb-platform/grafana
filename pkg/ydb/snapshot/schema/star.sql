CREATE TABLE IF NOT EXISTS `star` (
    `id` Serial8 NOT NULL,
    `user_id` Int64,
    `dashboard_id` Int64,
    `dashboard_uid` Text,
    `org_id` Int64,
    `updated` Datetime64,
    INDEX `UQE_star_user_id_dashboard_id` GLOBAL UNIQUE SYNC ON (`user_id`, `dashboard_id`),
    INDEX `UQE_star_user_id_dashboard_uid_org_id` GLOBAL UNIQUE SYNC ON (`user_id`, `dashboard_uid`, `org_id`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);


