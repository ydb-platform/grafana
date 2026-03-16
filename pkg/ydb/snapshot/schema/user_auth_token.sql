CREATE TABLE IF NOT EXISTS `user_auth_token` (
    `id` Serial8 NOT NULL,
    `user_id` Int64,
    `auth_token` Text,
    `prev_auth_token` Text,
    `user_agent` Text,
    `client_ip` Text,
    `auth_token_seen` Bool,
    `seen_at` Int64,
    `rotated_at` Int64,
    `created_at` Int64,
    `updated_at` Int64,
    `revoked_at` Int64,
    `external_session_id` Int64,
    INDEX `IDX_user_auth_token_revoked_at` GLOBAL ASYNC ON (`revoked_at`),
    INDEX `IDX_user_auth_token_user_id` GLOBAL ASYNC ON (`user_id`),
    INDEX `UQE_user_auth_token_auth_token` GLOBAL UNIQUE SYNC ON (`auth_token`),
    INDEX `UQE_user_auth_token_prev_auth_token` GLOBAL UNIQUE SYNC ON (`prev_auth_token`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);




