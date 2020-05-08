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
  `host_name` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '主机名',
  `host_UUID` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '主机 UUID',
  `os` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '操作系统',
  `platform` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '操作系统平台',
  `platform_family` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '操作系统家族',
  `platform_version` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '操作系统版本',
  `kernel_version` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '内核版本',
  `virtualization_system` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '虚拟化系统',
  `virtualization_role` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '虚拟化角色',
  `interfaces_json` varchar(2000) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '接口 json',
  `create_time` INTEGER DEFAULT NULL COMMENT '创建时间',
  `update_time` INTEGER DEFAULT NULL COMMENT '更新时间',
  `last_heartbeat_time` INTEGER DEFAULT NULL COMMENT '最后心跳时间',
  `config_version` INTEGER DEFAULT NULL COMMENT '配置版本号',
  PRIMARY KEY (`agent_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

DROP TABLE IF EXISTS `job`;
CREATE TABLE `job` (
  `job_id` int(11) NOT NULL AUTO_INCREMENT COMMENT '任务 ID',
  `agent_id` int(11) DEFAULT NULL,
  `runner_type` int(11) DEFAULT NULL COMMENT 'runner 类型： Telegraf = 0 DataDog = 1',
  `runner_module` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL COMMENT 'runner 模块',
  `module_version` varchar(100) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '模块版本',
  `running_interval` int(11) DEFAULT NULL COMMENT '运行间隔',
  `config` text COLLATE utf8mb4_bin COMMENT '配置',
  `create_time` int(11) DEFAULT NULL COMMENT '创建时间',
  `update_time` int(11) DEFAULT NULL COMMENT '更新时间',
  `running_state` int(11) DEFAULT NULL COMMENT '运行状态 Wait = 0 ,Ok = 1 ,Error = 2 ,Undef = 3',
  `error_msg` varchar(1000) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '错误信息',
  `last_running_time` int(11) DEFAULT NULL COMMENT '最近一次运行时间',
  `last_report_time` int(11) DEFAULT NULL COMMENT '运行状态上报时间',
  PRIMARY KEY (`job_id`),
  KEY `fk-job-host-on-agent_id` (`agent_id`),
  CONSTRAINT `fk-job-host-on-agent_id` FOREIGN KEY (`agent_id`) REFERENCES `host` (`agent_id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

SET FOREIGN_KEY_CHECKS = 1;
