CREATE TABLE IF NOT EXISTS `folder` (
  `id` Serial NOT NULL,
  `uid` Text NOT NULL,
  `org_id` Int64 NOT NULL,
  `title` Text NOT NULL,
  `description` Text,
  `parent_uid` Text,
  `created` Datetime64 NOT NULL,
  `updated` Datetime64 NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `UQE_folder_org_id_uid` GLOBAL UNIQUE SYNC ON (`org_id`, `uid`),
  INDEX `IDX_folder_org_id_parent_uid` GLOBAL SYNC ON (`org_id`, `parent_uid`)
);
