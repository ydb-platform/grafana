CREATE TABLE IF NOT EXISTS `resource_version` (
  `group` Text NOT NULL,
  `resource` Text NOT NULL,
  `resource_version` Int64 NOT NULL,
  PRIMARY KEY (`group`, `resource`),
);
