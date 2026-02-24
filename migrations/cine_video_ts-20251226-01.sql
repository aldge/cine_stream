-- +migrate Up
DROP TABLE IF EXISTS `cine_video_ts`;
CREATE TABLE `cine_video_ts` (
    `video_ts_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
    `video_id` char(32) NOT NULL DEFAULT '' COMMENT '视频id',
    `ts_sequence` int(10) unsigned NOT NULL DEFAULT '0' COMMENT 'TS序号',
    `ts_path` varchar(500) NOT NULL DEFAULT '' COMMENT 'TS文件存储路径',
    `duration` decimal(10,6) unsigned NOT NULL DEFAULT '0' COMMENT 'TS片段时长(秒)',
    `definition` varchar(50) NOT NULL DEFAULT '' COMMENT '清晰度',
    `create_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
    PRIMARY KEY(`video_ts_id`),
    KEY `video_id` (`video_id`),
    KEY `ts_sequence` (`ts_sequence`),
    KEY `definition` (`definition`),
    KEY `create_time` (`create_time`),
    UNIQUE KEY `video_id_sequence` (`video_id`, `ts_sequence`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='视频ts文件表';

-- +migrate Down
DROP TABLE IF EXISTS `cine_video_ts`;

