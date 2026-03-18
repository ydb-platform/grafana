CREATE TABLE IF NOT EXISTS `file` (
    `path` Text,
    `path_hash` Text,
    `parent_folder_path_hash` Text,
    `contents` Text,
    `etag` Text,
    `cache_control` Text,
    `content_disposition` Text,
    `updated` Datetime64,
    `created` Datetime64,
    `size` Int64,
    `mime_type` Text,
    INDEX `IDX_file_parent_folder_path_hash` GLOBAL ASYNC ON (`parent_folder_path_hash`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`path_hash`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);

