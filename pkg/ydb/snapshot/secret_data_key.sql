CREATE TABLE IF NOT EXISTS `secret_data_key` (
  `uid` Text NOT NULL,
  PRIMARY KEY (`uid`),
  `namespace` Text NOT NULL,
  `label` Text NOT NULL,
  `active` Int64 NOT NULL,
  `provider` Text NOT NULL,
  `encrypted_data` Bytes NOT NULL,
  `created` Datetime64 NOT NULL,
  `updated` Datetime64 NOT NULL,
  INDEX `IDX_secret_data_key_namespace_label_active` GLOBAL SYNC ON (`namespace`, `label`, `active`)
);
