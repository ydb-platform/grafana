CREATE TABLE IF NOT EXISTS `library_element_connection` (
  `id` Serial NOT NULL,
  `element_id` Int64 NOT NULL,
  `kind` Int64 NOT NULL,
  `connection_id` Int64 NOT NULL,
  `created` Datetime64 NOT NULL,
  `created_by` Int64 NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `UQE_library_element_connection_element_id_kind_connection_id` GLOBAL UNIQUE SYNC ON (`element_id`, `kind`, `connection_id`)
);
