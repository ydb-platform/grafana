CREATE TABLE IF NOT EXISTS `secrets` (
  `id` Serial NOT NULL,
  `org_id` Int64 NOT NULL,
  `namespace` Text NOT NULL,
  `type` Text NOT NULL,
  `value` Text,
  `created` Datetime64 NOT NULL,
  `updated` Datetime64 NOT NULL,
  PRIMARY KEY (`id`)
);
