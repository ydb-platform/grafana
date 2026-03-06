CREATE TABLE IF NOT EXISTS `query_history_star` (
  `id` Serial NOT NULL,
  `query_uid` Text NOT NULL,
  `user_id` Int64 NOT NULL,
  `org_id` Int64 NOT NULL DEFAULT 1,
  PRIMARY KEY (`id`),
  INDEX `UQE_query_history_star_user_id_query_uid` GLOBAL UNIQUE SYNC ON (`user_id`, `query_uid`)
);
