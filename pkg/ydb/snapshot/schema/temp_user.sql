CREATE TABLE IF NOT EXISTS `temp_user` (
    `id` Serial8 NOT NULL,
    `org_id` Int64,
    `version` Int64,
    `email` Text,
    `name` Text,
    `role` Text,
    `code` Text,
    `status` Text,
    `invited_by_user_id` Int64,
    `email_sent` Bool,
    `email_sent_on` Datetime64,
    `remote_addr` Text,
    `created` Int64,
    `updated` Int64,
    INDEX `IDX_temp_user_code` GLOBAL ASYNC ON (`code`),
    INDEX `IDX_temp_user_email` GLOBAL ASYNC ON (`email`),
    INDEX `IDX_temp_user_org_id` GLOBAL ASYNC ON (`org_id`),
    INDEX `IDX_temp_user_status` GLOBAL ASYNC ON (`status`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);




