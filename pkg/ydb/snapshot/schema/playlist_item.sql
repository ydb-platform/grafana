CREATE TABLE IF NOT EXISTS `playlist_item` (
  `id` Serial NOT NULL,
  `playlist_id` Int64 NOT NULL,
  `type` Text NOT NULL,
  `value` Text NOT NULL,
  `title` Text NOT NULL,
  `order` Int64 NOT NULL,
  PRIMARY KEY (`id`)
);
