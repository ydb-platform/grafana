CREATE TABLE IF NOT EXISTS `session` (
  `key` Text NOT NULL,
  PRIMARY KEY (`key`),
  `data` Bytes NOT NULL,
  `expiry` Int64 NOT NULL
);
