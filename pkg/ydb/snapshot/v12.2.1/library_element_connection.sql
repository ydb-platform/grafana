CREATE TABLE IF NOT EXISTS `library_element_connection` (
    `id` Serial8 NOT NULL,
    `element_id` Int64,
    `kind` Int64,
    `connection_id` Int64,
    `created` Datetime64,
    `created_by` Int64,
    INDEX `UQE_library_element_connection_element_id_kind_connection_id` GLOBAL UNIQUE SYNC ON (`element_id`, `kind`, `connection_id`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);

