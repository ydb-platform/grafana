CREATE TABLE IF NOT EXISTS `server_lock` (
  `id` Serial NOT NULL,
  `operation_uid` Text NOT NULL,
  `version` Int64 NOT NULL,
  `last_execution` Int64 NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `UQE_server_lock_operation_uid` GLOBAL UNIQUE SYNC ON (`operation_uid`)
);
