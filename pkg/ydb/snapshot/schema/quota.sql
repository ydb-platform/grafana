CREATE TABLE IF NOT EXISTS `quota` (
  `id` Serial NOT NULL,
  `org_id` Int64,
  `user_id` Int64,
  `target` Text NOT NULL,
  `limit` Int64 NOT NULL,
  `created` Datetime64 NOT NULL,
  `updated` Datetime64 NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `UQE_quota_org_id_user_id_target` GLOBAL UNIQUE SYNC ON (`org_id`, `user_id`, `target`)
);
