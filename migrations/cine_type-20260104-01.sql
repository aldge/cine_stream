-- +migrate Up
DROP TABLE IF EXISTS `cine_type`;
CREATE TABLE `cine_type` (
    `type_id` smallint unsigned NOT NULL AUTO_INCREMENT,
    `type_name` varchar(60) NOT NULL DEFAULT '',
    `type_en` varchar(60) NOT NULL DEFAULT '',
    `type_pid` smallint unsigned NOT NULL DEFAULT '0',
    `type_status` tinyint unsigned NOT NULL DEFAULT '1',
    PRIMARY KEY (`type_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8mb3;

-- +migrate Down
DROP TABLE IF EXISTS `cine_type`;