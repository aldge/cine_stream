-- ----------------------------------------------------------
-- 视频系统数据库迁移文件
-- 版本: 001
-- 创建时间: 2025-12-06
-- 描述: 创建视频相关表结构
-- ----------------------------------------------------------

CREATE DATABASE `cine_stream` DEFAULT CHARSET=utf8mb4;


-- ----------------------------------------------------------
-- 视频加密信息表（每个视频切片一个加密信息）
-- ----------------------------------------------------------
DROP TABLE IF EXISTS `cine_video_encrypt`;
CREATE TABLE `cine_video_encrypt` (
	`video_encrypt_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '主键id',
	`video_id` char(32) NOT NULL DEFAULT '' COMMENT '视频id',
	`key` char(64) NOT NULL DEFAULT '' COMMENT '加密 key',
	`iv` char(64) NOT NULL DEFAULT '' COMMENT '加密向量',
	`create_time` int(11) unsigned NOT NULL DEFAULT '0' COMMENT '创建时间',
	PRIMARY KEY(`video_encrypt_id`),
	KEY `video_id` (`video_id`),
	KEY `create_time` (`create_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='视频加密信息表';

-- ----------------------------------------------------------
-- 视频 ts 文件表(vod_id分表)
-- ----------------------------------------------------------
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