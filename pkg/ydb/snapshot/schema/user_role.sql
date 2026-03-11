CREATE TABLE IF NOT EXISTS `user_role` (
  `id` Serial NOT NULL,
  `org_id` Int64 NOT NULL,
  `user_id` Int64 NOT NULL,
  `role_id` Int64 NOT NULL,
  `created` Datetime64 NOT NULL,
  `group_mapping_uid` Text DEFAULT '',
  PRIMARY KEY (`id`),
  INDEX `IDX_user_role_org_id` GLOBAL SYNC ON (`org_id`),
  INDEX `IDX_user_role_user_id` GLOBAL SYNC ON (`user_id`),
  INDEX `UQE_user_role_org_id_user_id_role_id_group_mapping_uid` GLOBAL UNIQUE SYNC ON (
    `org_id`,
    `user_id`,
    `role_id`,
    `group_mapping_uid`
  ),
);
