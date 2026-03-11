CREATE TABLE IF NOT EXISTS `secret_migration_log` (
  `id` Serial NOT NULL,
  `migration_id` Text NOT NULL,
  `sql` Text NOT NULL,
  `success` Int64 NOT NULL,
  `error` Text NOT NULL,
  `timestamp` Datetime64 NOT NULL,
  PRIMARY KEY (`id`)
);
