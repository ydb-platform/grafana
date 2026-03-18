CREATE TABLE IF NOT EXISTS `team` (
    `id` Serial8 NOT NULL,
    `name` Text,
    `org_id` Int64,
    `created` Datetime64,
    `updated` Datetime64,
    `email` Text,
    `uid` Text,
    `external_uid` Text,
    `is_provisioned` Bool,
    INDEX `UQE_team_org_id_uid` GLOBAL UNIQUE SYNC ON (`org_id`, `uid`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);

