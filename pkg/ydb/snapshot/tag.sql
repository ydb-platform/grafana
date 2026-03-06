CREATE TABLE IF NOT EXISTS `tag` (
  `id` Serial NOT NULL,
  `key` Text NOT NULL,
  `value` Text NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `UQE_tag_key_value` GLOBAL UNIQUE SYNC ON (`key`, `value`)
);
