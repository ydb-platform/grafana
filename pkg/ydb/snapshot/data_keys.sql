CREATE TABLE IF NOT EXISTS `data_keys` (
  `name` Text NOT NULL,
  `active` Int64 NOT NULL,
  `scope` Text NOT NULL,
  `provider` Text NOT NULL,
  `encrypted_data` Bytes NOT NULL,
  `created` Datetime64 NOT NULL,
  `updated` Datetime64 NOT NULL,
  `label` Text NOT NULL DEFAULT '',
  PRIMARY KEY (`name`)
);
