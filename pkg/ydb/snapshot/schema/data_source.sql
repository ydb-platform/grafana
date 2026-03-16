CREATE TABLE IF NOT EXISTS `data_source` (
    `id` Serial8 NOT NULL,
    `org_id` Int64,
    `version` Int64,
    `type` Text,
    `name` Text,
    `access` Text,
    `url` Text,
    `password` Text,
    `user` Text,
    `database` Text,
    `basic_auth` Bool,
    `basic_auth_user` Text,
    `basic_auth_password` Text,
    `is_default` Bool,
    `json_data` Text,
    `created` Datetime64,
    `updated` Datetime64,
    `with_credentials` Bool,
    `secure_json_data` Text,
    `read_only` Bool,
    `uid` Text,
    `is_prunable` Bool,
    `api_version` Text,
    INDEX `IDX_data_source_org_id` GLOBAL ASYNC ON (`org_id`),
    INDEX `IDX_data_source_org_id_is_default` GLOBAL ASYNC ON (`org_id`, `is_default`),
    INDEX `UQE_data_source_org_id_name` GLOBAL UNIQUE SYNC ON (`org_id`, `name`),
    INDEX `UQE_data_source_org_id_uid` GLOBAL UNIQUE SYNC ON (`org_id`, `uid`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);




