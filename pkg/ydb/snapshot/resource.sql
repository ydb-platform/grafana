CREATE TABLE IF NOT EXISTS `resource` (
  `guid` Text NOT NULL,
  `resource_version` Int64,
  `group` Text NOT NULL,
  `resource` Text NOT NULL,
  `namespace` Text NOT NULL,
  `name` Text NOT NULL,
  `value` Text,
  `action` Int64 NOT NULL,
  `label_set` Text,
  `previous_resource_version` Int64,
  `folder` Text NOT NULL DEFAULT '',
  PRIMARY KEY (`guid`),
  INDEX `UQE_resource_namespace_group_resource_name` GLOBAL UNIQUE SYNC ON (`namespace`, `group`, `resource`, `name`),
  INDEX `UQE_resource_history_namespace_group_name_version` GLOBAL UNIQUE SYNC ON (
    `namespace`,
    `group`,
    `resource`,
    `name`,
    `resource_version`
  ),
  INDEX `UQE_resource_version_group_resource` GLOBAL UNIQUE SYNC ON (`group`, `resource`),
  INDEX `IDX_resource_history_namespace_group_name` GLOBAL SYNC ON (`namespace`, `group`, `resource`, `name`),
  INDEX `IDX_resource_history_group_resource_resource_version` GLOBAL SYNC ON (`group`, `resource`, `resource_version`),
  INDEX `IDX_resource_group_resource` GLOBAL SYNC ON (`group`, `resource`),
  INDEX `IDX_resource_history_namespace_group_resource_action_version` GLOBAL SYNC ON (
    `namespace`,
    `group`,
    `resource`,
    `action`,
    `resource_version`
  ),
);
