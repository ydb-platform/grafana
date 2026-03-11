CREATE TABLE IF NOT EXISTS `anon_device` (
  `id` Serial NOT NULL,
  `client_ip` Text NOT NULL,
  `created_at` Datetime64 NOT NULL,
  `device_id` Text NOT NULL,
  `updated_at` Datetime64 NOT NULL,
  `user_agent` Text NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `UQE_anon_device_device_id` GLOBAL UNIQUE SYNC ON (`device_id`),
  INDEX `IDX_anon_device_updated_at` GLOBAL SYNC ON (`updated_at`)
);
