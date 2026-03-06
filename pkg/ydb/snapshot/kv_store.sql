CREATE TABLE IF NOT EXISTS `kv_store` (
  `id` Serial NOT NULL,
  `org_id` Int64 NOT NULL,
  `namespace` Text NOT NULL,
  `key` Text NOT NULL,
  `value` Text NOT NULL,
  `created` Datetime64 NOT NULL,
  `updated` Datetime64 NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `UQE_kv_store_org_id_namespace_key` GLOBAL UNIQUE SYNC ON (`org_id`, `namespace`, `key`)
);
