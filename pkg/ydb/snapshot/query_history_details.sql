CREATE TABLE IF NOT EXISTS `query_history_details` (
  `id` Serial NOT NULL,
  `query_history_item_uid` Text NOT NULL,
  `datasource_uid` Text NOT NULL,
  PRIMARY KEY (`id`)
);
