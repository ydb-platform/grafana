CREATE TABLE IF NOT EXISTS `alert_rule_tag` (
  `id` Serial NOT NULL,
  `alert_id` Int64 NOT NULL,
  `tag_id` Int64 NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `UQE_alert_rule_tag_alert_id_tag_id` GLOBAL UNIQUE SYNC ON (`alert_id`, `tag_id`),
  INDEX `IDX_alert_rule_tag_alert_id` GLOBAL SYNC ON (`alert_id`)
);
