CREATE TABLE IF NOT EXISTS `resource_blob` (
  `uuid` UUID NOT NULL,
  `created` Datetime64 NOT NULL,
  `group` Text NOT NULL,
  `resource` Text NOT NULL,
  `namespace` Text NOT NULL,
  `name` Text NOT NULL,
  `value` Bytes NOT NULL,
  `hash` Text NOT NULL,
  `content_type` Text NOT NULL,
  PRIMARY KEY (`uuid`),
  INDEX `IDX_resource_history_namespace_group_name` GLOBAL SYNC ON (`namespace`, `group`, `resource`, `name`),
  INDEX `IDX_resource_blob_created` GLOBAL SYNC ON (`created`)
);
