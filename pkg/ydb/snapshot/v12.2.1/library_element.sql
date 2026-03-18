CREATE TABLE IF NOT EXISTS `library_element` (
    `id` Serial8 NOT NULL,
    `org_id` Int64,
    `folder_id` Int64,
    `uid` Text,
    `name` Text,
    `kind` Int64,
    `type` Text,
    `description` Text,
    `model` Text,
    `created` Datetime64,
    `created_by` Int64,
    `updated` Datetime64,
    `updated_by` Int64,
    `version` Int64,
    `folder_uid` Text,
    INDEX `UQE_library_element_org_id_folder_id_name_kind` GLOBAL UNIQUE SYNC ON (`org_id`, `folder_id`, `name`, `kind`),
    INDEX `UQE_library_element_org_id_folder_uid_name_kind` GLOBAL UNIQUE SYNC ON (`org_id`, `folder_uid`, `name`, `kind`),
    INDEX `UQE_library_element_org_id_uid` GLOBAL UNIQUE SYNC ON (`org_id`, `uid`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);



