CREATE TABLE IF NOT EXISTS `org_user` (
  `id` Serial NOT NULL,
  `org_id` Int64 NOT NULL,
  `user_id` Int64 NOT NULL,
  `role` Text NOT NULL,
  `created` Datetime64 NOT NULL,
  `updated` Datetime64 NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `IDX_org_user_org_id` GLOBAL SYNC ON (`org_id`),
  INDEX `UQE_org_user_org_id_user_id` GLOBAL UNIQUE SYNC ON (`org_id`, `user_id`),
  INDEX `IDX_org_user_user_id` GLOBAL SYNC ON (`user_id`),
);
