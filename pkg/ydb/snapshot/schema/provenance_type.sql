CREATE TABLE IF NOT EXISTS `provenance_type` (
    `id` Serial8 NOT NULL,
    `org_id` Int64,
    `record_key` Text,
    `record_type` Text,
    `provenance` Text,
    INDEX `UQE_provenance_type_record_type_record_key_org_id` GLOBAL UNIQUE SYNC ON (`record_type`, `record_key`, `org_id`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);

