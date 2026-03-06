CREATE TABLE IF NOT EXISTS `signing_key` (
  `id` Serial NOT NULL,
  `key_id` Text NOT NULL,
  `private_key` Text NOT NULL,
  `added_at` Datetime64 NOT NULL,
  `expires_at` Datetime64,
  `alg` Text NOT NULL,
  PRIMARY KEY (`id`),
  INDEX `UQE_signing_key_key_id` GLOBAL UNIQUE SYNC ON (`key_id`)
);
