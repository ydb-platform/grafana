CREATE TABLE IF NOT EXISTS `user` (
    `id` Serial8 NOT NULL,
    `version` Int64,
    `login` Text,
    `email` Text,
    `name` Text,
    `password` Text,
    `salt` Text,
    `rands` Text,
    `company` Text,
    `org_id` Int64,
    `is_admin` Bool,
    `email_verified` Bool,
    `theme` Text,
    `created` Datetime64,
    `updated` Datetime64,
    `help_flags1` Int64,
    `last_seen_at` Datetime64,
    `is_disabled` Bool,
    `is_service_account` Bool,
    `uid` Text,
    `is_provisioned` Bool,
    INDEX `IDX_user_is_service_account_last_seen_at` GLOBAL ASYNC ON (`is_service_account`, `last_seen_at`),
    INDEX `IDX_user_login_email` GLOBAL ASYNC ON (`login`, `email`),
    INDEX `UQE_user_email` GLOBAL UNIQUE SYNC ON (`email`),
    INDEX `UQE_user_login` GLOBAL UNIQUE SYNC ON (`login`),
    INDEX `UQE_user_uid` GLOBAL UNIQUE SYNC ON (`uid`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);





