CREATE TABLE IF NOT EXISTS `login_attempt` (
  `id` Serial NOT NULL,
  `username` Text NOT NULL,
  `ip_address` Text NOT NULL,
  `created` Int64 NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  INDEX `IDX_login_attempt_username` GLOBAL SYNC ON (`username`)
);
