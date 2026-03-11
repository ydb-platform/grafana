CREATE TABLE IF NOT EXISTS `annotation_tag` (
  `id` Serial NOT NULL,
  `annotation_id` Int64 NOT NULL,
  `tag_id` Int64 NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `UQE_annotation_tag_annotation_id_tag_id` GLOBAL UNIQUE SYNC ON (`annotation_id`, `tag_id`)
);
