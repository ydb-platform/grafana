CREATE TABLE IF NOT EXISTS `user_external_session` (
    `id` Serial8 NOT NULL,
    `user_auth_id` Int64,
    `user_id` Int64,
    `auth_module` Text,
    `access_token` Text,
    `id_token` Text,
    `refresh_token` Text,
    `session_id` Text,
    `session_id_hash` Text,
    `name_id` Text,
    `name_id_hash` Text,
    `expires_at` Datetime64,
    `created_at` Datetime64,
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);
