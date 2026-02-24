-- +migrate Up
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

-- +migrate Down
DROP TABLE IF EXISTS `cine_video_encrypt`;

