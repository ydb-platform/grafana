CREATE TABLE IF NOT EXISTS `alert_image` (
  `id` Serial NOT NULL,
  `token` Text NOT NULL,
  `path` Text NOT NULL,
  `url` Text NOT NULL,
  `created_at` Datetime64 NOT NULL,
  `expires_at` Datetime64 NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `UQE_alert_image_token` GLOBAL UNIQUE SYNC ON (`token`)
);
