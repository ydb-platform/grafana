CREATE TABLE IF NOT EXISTS `star` (
  `id` Serial NOT NULL,
  `user_id` Int64 NOT NULL,
  `dashboard_id` Int64 NOT NULL,
  `dashboard_uid` Text,
  `org_id` Int64 DEFAULT 1,
  `updated` Datetime64,
  PRIMARY KEY (`id`),
  INDEX `UQE_star_user_id_dashboard_id` GLOBAL UNIQUE SYNC ON (`user_id`, `dashboard_id`),
  INDEX `UQE_star_user_id_dashboard_uid_org_id` GLOBAL UNIQUE SYNC ON (`user_id`, `dashboard_uid`, `org_id`)
);
