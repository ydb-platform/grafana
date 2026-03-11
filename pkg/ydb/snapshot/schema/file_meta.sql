CREATE TABLE IF NOT EXISTS `file_meta` (
  `path_hash` Text NOT NULL,
  `key` Text NOT NULL,
  `value` Text NOT NULL,
  PRIMARY KEY (`path_hash`, `key`)
);
