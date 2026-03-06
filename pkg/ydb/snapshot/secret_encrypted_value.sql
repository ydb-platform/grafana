CREATE TABLE IF NOT EXISTS `secret_encrypted_value` (
  `namespace` Text NOT NULL,
  `name` Text NOT NULL,
  `version` Int64 NOT NULL,
  `encrypted_data` Bytes NOT NULL,
  `created` Int64 NOT NULL,
  `updated` Int64 NOT NULL,
  PRIMARY KEY (`namespace`, `name`, `version`)
);
