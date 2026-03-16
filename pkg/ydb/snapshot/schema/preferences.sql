CREATE TABLE IF NOT EXISTS `preferences` (
    `id` Serial8 NOT NULL,
    `org_id` Int64,
    `user_id` Int64,
    `version` Int64,
    `home_dashboard_id` Int64,
    `timezone` Text,
    `theme` Text,
    `created` Datetime64,
    `updated` Datetime64,
    `team_id` Int64,
    `week_start` Text,
    `json_data` Text,
    `home_dashboard_uid` Text,
    INDEX `IDX_preferences_org_id` GLOBAL ASYNC ON (`org_id`),
    INDEX `IDX_preferences_user_id` GLOBAL ASYNC ON (`user_id`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);


