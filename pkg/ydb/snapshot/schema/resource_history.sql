CREATE TABLE IF NOT EXISTS `resource_history` (
  `guid` Text NOT NULL,
  PRIMARY KEY (`guid`),
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
  `generation` Int64 NOT NULL DEFAULT 0,
  INDEX `UQE_resource_history_namespace_group_name_version` GLOBAL UNIQUE SYNC ON (
    `namespace`,
    `group`,
    `resource`,
    `name`,
    `resource_version`
  ),
  INDEX `IDX_resource_history_resource_version` GLOBAL SYNC ON (`resource_version`),
  INDEX `IDX_resource_history_group_resource_resource_version` GLOBAL SYNC ON (`group`, `resource`, `resource_version`),
  INDEX `IDX_resource_history_namespace_group_resource_action_version` GLOBAL SYNC ON (
    `namespace`,
    `group`,
    `resource`,
    `action`,
    `resource_version`
  ),
  INDEX `IDX_resource_history_namespace_group_resource_name_generation` GLOBAL SYNC ON (
    `namespace`,
    `group`,
    `resource`,
    `name`,
    `generation`
  ),
);
