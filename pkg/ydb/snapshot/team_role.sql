CREATE TABLE IF NOT EXISTS `team_role` (
  `id` Serial NOT NULL,
  `org_id` Int64 NOT NULL,
  `team_id` Int64 NOT NULL,
  `role_id` Int64 NOT NULL,
  `created` Datetime64 NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `IDX_team_role_org_id` GLOBAL SYNC ON (`org_id`),
  INDEX `UQE_team_role_org_id_team_id_role_id` GLOBAL UNIQUE SYNC ON (`org_id`, `team_id`, `role_id`),
  INDEX `IDX_team_role_team_id` GLOBAL SYNC ON (`team_id`),
);
