CREATE TABLE IF NOT EXISTS `cache_data` (
  `cache_key` Text NOT NULL,
  `data` Bytes NOT NULL,
  `expires` Int64 NOT NULL,
  `created_at` Int64 NOT NULL,
  PRIMARY KEY (`cache_key`),
);
