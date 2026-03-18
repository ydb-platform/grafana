CREATE TABLE IF NOT EXISTS `user_auth` (
    `id` Serial8 NOT NULL,
    `user_id` Int64,
    `auth_module` Text,
    `auth_id` Text,
    `created` Datetime64,
    `o_auth_access_token` Text,
    `o_auth_refresh_token` Text,
    `o_auth_token_type` Text,
    `o_auth_expiry` Datetime64,
    `o_auth_id_token` Text,
    `external_uid` Text,
    INDEX `IDX_user_auth_auth_module_auth_id` GLOBAL ASYNC ON (`auth_module`, `auth_id`),
    INDEX `IDX_user_auth_user_id` GLOBAL ASYNC ON (`user_id`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);


