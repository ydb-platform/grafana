CREATE TABLE IF NOT EXISTS `preferences` (
  `id` Serial NOT NULL,
  `org_id` Int64 NOT NULL,
  `user_id` Int64 NOT NULL,
  `version` Int64 NOT NULL,
  `home_dashboard_id` Int64 NOT NULL,
  `timezone` Text NOT NULL,
  `theme` Text NOT NULL,
  `created` Datetime64 NOT NULL,
  `updated` Datetime64 NOT NULL,
  `team_id` Int64,
  `week_start` Text,
  `json_data` Text,
  `home_dashboard_uid` Text,
  PRIMARY KEY (`id`),
  INDEX `IDX_preferences_org_id` GLOBAL SYNC ON (`org_id`),
  INDEX `IDX_preferences_user_id` GLOBAL SYNC ON (`user_id`)
);
