CREATE TABLE IF NOT EXISTS `cloud_migration_session` (
    `id` Serial8 NOT NULL,
    `uid` Text,
    `auth_token` Text,
    `slug` Text,
    `stack_id` Int64,
    `region_slug` Text,
    `cluster_slug` Text,
    `created` Datetime64,
    `updated` Datetime64,
    `org_id` Int64,
    INDEX `UQE_cloud_migration_session_uid` GLOBAL UNIQUE SYNC ON (`uid`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);

