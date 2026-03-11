CREATE TABLE IF NOT EXISTS `test_data` (
  `id` Serial NOT NULL,
  `metric1` Text,
  `metric2` Text,
  `value_big_int` Int64,
  `value_double` Double,
  `value_float` Double,
  `value_int` Int64,
  `time_epoch` Int64 NOT NULL,
  `time_date_time` Datetime64 NOT NULL,
  `time_time_stamp` Datetime64 NOT NULL,
  PRIMARY KEY (`id`)
);
