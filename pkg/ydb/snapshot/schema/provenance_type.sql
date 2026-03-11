CREATE TABLE IF NOT EXISTS `provenance_type` (
  `id` Serial NOT NULL,
  `org_id` Int64 NOT NULL,
  `record_key` Text NOT NULL,
  `record_type` Text NOT NULL,
  `provenance` Text NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `UQE_provenance_type_record_type_record_key_org_id` GLOBAL UNIQUE SYNC ON (`record_type`, `record_key`, `org_id`)
);
