CREATE TABLE IF NOT EXISTS `secret_secure_value` (
    `guid` Text,
    `name` Text,
    `namespace` Text,
    `annotations` Text,
    `labels` Text,
    `created` Int64,
    `created_by` Text,
    `updated` Int64,
    `updated_by` Text,
    `external_id` Text,
    `active` Bool,
    `version` Int64,
    `description` Text,
    `keeper` Text,
    `decrypters` Text,
    `ref` Text,
    `owner_reference_api_group` Text,
    `owner_reference_api_version` Text,
    `owner_reference_kind` Text,
    `owner_reference_name` Text,
    `lease_token` Text,
    `lease_created` Int64,
    INDEX `IDX_secret_secure_value_lease_created` GLOBAL ASYNC ON (`lease_created`),
    INDEX `IDX_secret_secure_value_lease_token` GLOBAL ASYNC ON (`lease_token`),
    INDEX `IDX_secret_secure_value_namespace_active_updated` GLOBAL ASYNC ON (`namespace`, `active`, `updated`),
    INDEX `UQE_secret_secure_value_namespace_name_version` GLOBAL UNIQUE SYNC ON (`namespace`, `name`, `version`),
    INDEX `UQE_secret_secure_value_namespace_name_version_active` GLOBAL UNIQUE SYNC ON (`namespace`, `name`, `version`, `active`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`guid`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);





