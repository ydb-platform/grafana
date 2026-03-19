CREATE TABLE IF NOT EXISTS `test_data` (
    `id` Serial8 NOT NULL,
    `metric1` Text,
    `metric2` Text,
    `value_big_int` Int64,
    `value_double` Double,
    `value_float` Double,
    `value_int` Int64,
    `time_epoch` Int64,
    `time_date_time` Datetime64,
    `time_time_stamp` Datetime64,
    FAMILY `default` (COMPRESSION = 'off'),
    PRIMARY KEY (`id`)
)
WITH (
    AUTO_PARTITIONING_BY_SIZE = ENABLED,
    AUTO_PARTITIONING_PARTITION_SIZE_MB = 2048,
    AUTO_PARTITIONING_BY_LOAD = ENABLED
);
