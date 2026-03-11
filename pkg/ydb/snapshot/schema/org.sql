CREATE TABLE IF NOT EXISTS `org` (
  `id` Int64 NOT NULL,
  `version` Int64 NOT NULL,
  `name` Text NOT NULL,
  `address1` Text,
  `address2` Text,
  `city` Text,
  `state` Text,
  `zip_code` Text,
  `country` Text,
  `billing_email` Text,
  `created` Datetime64 NOT NULL,
  `updated` Datetime64 NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `UQE_org_name` GLOBAL UNIQUE SYNC ON (`name`)
);
