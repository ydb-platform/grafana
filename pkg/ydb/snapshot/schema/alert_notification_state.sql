CREATE TABLE IF NOT EXISTS `alert_notification_state` (
  `id` Serial NOT NULL,
  `org_id` Int64 NOT NULL,
  `alert_id` Int64 NOT NULL,
  `notifier_id` Int64 NOT NULL,
  `state` Text NOT NULL,
  `version` Int64 NOT NULL,
  `updated_at` Int64 NOT NULL,
  `alert_rule_state_updated_version` Int64 NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `UQE_alert_notification_state_org_id_alert_id_notifier_id` GLOBAL UNIQUE SYNC ON (`org_id`, `alert_id`, `notifier_id`),
  INDEX `IDX_alert_notification_state_alert_id` GLOBAL SYNC ON (`alert_id`)
);
