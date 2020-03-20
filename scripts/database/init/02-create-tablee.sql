/*
 Navicat Premium Data Transfer

 Source Server         : gringotts-db-root
 Source Server Type    : MySQL
 Source Server Version : 50729
 Source Host           : 127.0.0.1:3306
 Source Schema         : gringotts

 Target Server Type    : MySQL
 Target Server Version : 50729
 File Encoding         : 65001

 Date: 18/03/2020 12:40:14
*/

USE gringotts;
SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for host
-- ----------------------------
DROP TABLE IF EXISTS `host`;
CREATE TABLE `host` (
  `agent_id` int(11) NOT NULL AUTO_INCREMENT COMMENT 'AgentID',
  `host_name` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '主机名',
  `host_UUID` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '主机 UUID',
  `os` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '操作系统',
  `platform` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '操作系统平台',
  `platform_family` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '操作系统家族',
  `platform_version` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '操作系统版本',
  `kernel_version` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '内核版本',
  `virtualization_system` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '虚拟化系统',
  `virtualization_role` varchar(100) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '虚拟化角色',
  `interfaces_json` varchar(2000) COLLATE utf8mb4_unicode_ci DEFAULT NULL COMMENT '接口 json',
  `create_time` INTEGER DEFAULT NULL COMMENT '创建时间',
  `update_time` INTEGER DEFAULT NULL COMMENT '更新时间',
  `last_heartbeat_time` INTEGER DEFAULT NULL COMMENT '最后心跳时间',
  PRIMARY KEY (`agent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

SET FOREIGN_KEY_CHECKS = 1;
