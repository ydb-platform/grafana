CREATE TABLE IF NOT EXISTS `entity_event` (
  `id` Serial NOT NULL,
  `entity_id` Text NOT NULL,
  `event_type` Text NOT NULL,
  `created` Int64 NOT NULL,
  PRIMARY KEY (`id`)
);
