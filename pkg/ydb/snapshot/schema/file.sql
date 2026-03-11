CREATE TABLE IF NOT EXISTS `file` (
  `path` Text NOT NULL,
  `path_hash` Text NOT NULL,
  `parent_folder_path_hash` Text NOT NULL,
  `contents` Bytes NOT NULL,
  `etag` Text NOT NULL,
  `cache_control` Text NOT NULL,
  `content_disposition` Text NOT NULL,
  `updated` Datetime64 NOT NULL,
  `created` Datetime64 NOT NULL,
  `size` Int64 NOT NULL,
  `mime_type` Text NOT NULL,
  PRIMARY KEY (`path_hash`),
  INDEX `IDX_file_parent_folder_path_hash` GLOBAL SYNC ON (`parent_folder_path_hash`)
);
