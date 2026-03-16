CREATE TABLE IF NOT EXISTS `plugin_setting` (
    `id` Serial8 NOT NULL,
    `org_id` Int64,
    `plugin_id` Text,
    `enabled` Bool,
    `pinned` Bool,
    `json_data` Text,
    `secure_json_data` Text,
    `created` Datetime64,
    `updated` Datetime64,
    `plugin_version` Text,
    INDEX `UQE_plugin_setting_org_id_plugin_id` GLOBAL UNIQUE SYNC ON (`org_id`, `plugin_id`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);

