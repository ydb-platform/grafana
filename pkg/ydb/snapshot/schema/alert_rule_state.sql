CREATE TABLE IF NOT EXISTS `alert_rule_state` (
  `id` Serial NOT NULL,
  `org_id` Int64 NOT NULL,
  `rule_uid` Text NOT NULL,
  `data` Bytes NOT NULL,
  `updated_at` Datetime64 NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `UQE_alert_rule_state_org_id_rule_uid` GLOBAL UNIQUE SYNC ON (`org_id`, `rule_uid`)
);
