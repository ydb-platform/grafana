CREATE TABLE IF NOT EXISTS `alert_instance` (
  `rule_org_id` Int64 NOT NULL,
  `rule_uid` Text NOT NULL,
  `labels` Text NOT NULL,
  `labels_hash` Text NOT NULL,
  `current_state` Text NOT NULL,
  `current_state_since` Int64 NOT NULL,
  `last_eval_time` Int64 NOT NULL,
  `current_state_end` Int64 NOT NULL DEFAULT 0,
  `current_reason` Text,
  `result_fingerprint` Text,
  `resolved_at` Int64,
  `last_sent_at` Int64,
  `fired_at` Int64,
  PRIMARY KEY (`rule_org_id`, `rule_uid`, `labels_hash`),
  INDEX `IDX_alert_instance_rule_org_id_rule_uid_current_state` GLOBAL SYNC ON (`rule_org_id`, `rule_uid`, `current_state`),
  INDEX `IDX_alert_instance_rule_org_id_current_state` GLOBAL SYNC ON (`rule_org_id`, `current_state`)
);
