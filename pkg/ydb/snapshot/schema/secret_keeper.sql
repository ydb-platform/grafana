CREATE TABLE IF NOT EXISTS `secret_keeper` (
  `guid` Text NOT NULL,
  PRIMARY KEY (`guid`),
  `name` Text NOT NULL,
  `namespace` Text NOT NULL,
  `annotations` Text,
  `labels` Text,
  `created` Int64 NOT NULL,
  `created_by` Text NOT NULL,
  `updated` Int64 NOT NULL,
  `updated_by` Text NOT NULL,
  `description` Text NOT NULL,
  `type` Text NOT NULL,
  `payload` Text,
  INDEX `UQE_secret_keeper_namespace_name` GLOBAL UNIQUE SYNC ON (`namespace`, `name`)
);
