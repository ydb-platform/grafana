CREATE TABLE IF NOT EXISTS `org` (
    `id` Serial8 NOT NULL,
    `version` Int64,
    `name` Text,
    `address1` Text,
    `address2` Text,
    `city` Text,
    `state` Text,
    `zip_code` Text,
    `country` Text,
    `billing_email` Text,
    `created` Datetime64,
    `updated` Datetime64,
    INDEX `UQE_org_name` GLOBAL UNIQUE SYNC ON (`name`),
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);

