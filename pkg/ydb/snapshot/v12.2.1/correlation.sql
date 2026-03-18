CREATE TABLE IF NOT EXISTS `correlation` (
    `uid` Text,
    `org_id` Int64,
    `source_uid` Text,
    `target_uid` Text,
    `label` Text,
    `description` Text,
    `config` Text,
    `provisioned` Bool,
    `type` Text,
    INDEX `IDX_correlation_org_id` GLOBAL ASYNC ON (`org_id`),
    INDEX `IDX_correlation_source_uid` GLOBAL ASYNC ON (`source_uid`),
    INDEX `IDX_correlation_uid` GLOBAL ASYNC ON (`uid`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`uid`, `org_id`, `source_uid`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);



