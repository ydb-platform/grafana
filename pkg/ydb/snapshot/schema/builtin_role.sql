CREATE TABLE IF NOT EXISTS `builtin_role` (
  `id` Serial NOT NULL,
  `role` Text NOT NULL,
  `role_id` Int64 NOT NULL,
  `created` Datetime64 NOT NULL,
  `updated` Datetime64 NOT NULL,
  `org_id` Int64 NOT NULL DEFAULT 0,
  PRIMARY KEY (`id`),
  INDEX `IDX_builtin_role_role_id` GLOBAL SYNC ON (`role_id`),
  INDEX `IDX_builtin_role_role` GLOBAL SYNC ON (`role`),
  INDEX `IDX_builtin_role_org_id` GLOBAL SYNC ON (`org_id`),
  INDEX `UQE_builtin_role_org_id_role_id_role` GLOBAL UNIQUE SYNC ON (`org_id`, `role_id`, `role`),
);
