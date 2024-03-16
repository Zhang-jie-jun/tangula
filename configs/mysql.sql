drop database if exists `tangula`;
create database `tangula` default character set utf8mb4 collate utf8mb4_unicode_ci;

use `tangula`;

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS=0;


-- ----------------------------
-- Table structure for tangula_role
-- ----------------------------
DROP TABLE IF EXISTS `tangula_role`;
CREATE TABLE `tangula_role` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '名称',
  `desc` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci COMMENT '描述',
  `created_time` datetime NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '角色信息';

-- ----------------------------
-- Table structure for tangula_user
-- ----------------------------
DROP TABLE IF EXISTS `tangula_user`;
CREATE TABLE `tangula_user` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '名称',
  `account` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '账号',
  `mail` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '邮箱地址',
  `phone` varchar(11) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '联系方式',
  `status` int DEFAULT 1 NOT NULL COMMENT '角色状态',
  `usage_count` int DEFAULT 0 COMMENT '访问计数',
  `role_id` int(11) NOT NULL COMMENT '角色类型ID',
  `created_time` datetime NOT NULL COMMENT '创建时间',
  `updated_time` datetime NOT NULL COMMENT '更新时间',
  PRIMARY KEY (`id`),
  KEY `role_id` (`role_id`),
  CONSTRAINT `tangula_user_ibfk_1` FOREIGN KEY (`role_id`) REFERENCES `tangula_role` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '用户信息';

-- ----------------------------
-- Table structure for tangula_host
-- ----------------------------
DROP TABLE IF EXISTS `tangula_host`;
CREATE TABLE `tangula_host` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '名称',
  `desc` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '描述',
  `hostname` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主机名',
  `type` int NOT NULL COMMENT '主机类型',
  `os` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '操作系统类型',
  `arch` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '系统架构',
  `ip` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主机IP地址',
  `port` int DEFAULT 22 NOT NULL COMMENT '主机端口',
  `username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主机登录用户名',
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主机登录密码',
  `create_user` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '主机创建用户',
  `auth_type` int DEFAULT 111 NOT NULL COMMENT '主机权限类型',
  `created_time` datetime NOT NULL COMMENT '创建时间',
  `updated_time` datetime NOT NULL COMMENT '更新时间',
  `deleted_time` datetime COMMENT '删除时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '主机信息';

-- ----------------------------
-- Table structure for tangula_platform
-- ----------------------------
DROP TABLE IF EXISTS `tangula_platform`;
CREATE TABLE `tangula_platform` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '名称',
  `desc` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '描述',
  `type` int NOT NULL COMMENT '平台类型',
  `ip` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '平台IP地址',
  `port` int DEFAULT 22 NOT NULL COMMENT '平台端口',
  `username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '平台登录用户名',
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '平台登录密码',
  `version` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '平台版本信息',
  `create_user` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '平台创建用户',
  `auth_type` int DEFAULT 111 NOT NULL COMMENT '主机权限类型',
  `created_time` datetime NOT NULL COMMENT '创建时间',
  `updated_time` datetime NOT NULL COMMENT '更新时间',
  `deleted_time` datetime COMMENT '删除时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '平台信息';

-- ----------------------------
-- Table structure for tangula_tenant
-- ----------------------------
DROP TABLE IF EXISTS `tangula_tenant`;
CREATE TABLE `tangula_tenant` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '名称',
  `desc` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '描述',
  `domain_id` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '可用域ID',
  `username` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '平台登录用户名',
  `password` varchar(255) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '平台登录密码',
  `platform_id` int(11) NOT NULL COMMENT '平台ID',
  PRIMARY KEY (`id`),
  KEY `platform_id` (`platform_id`),
  CONSTRAINT `tangula_platform_ibfk_1` FOREIGN KEY (`platform_id`) REFERENCES `tangula_platform` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '租户信息';


-- ----------------------------
-- Table structure for tangula_record
-- ----------------------------
DROP TABLE IF EXISTS `tangula_record`;
CREATE TABLE `tangula_record` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `operation` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '操作内容',
  `object` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT 'unknown' COMMENT '操作对象',
  `detail` varchar(4096) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '操作详情',
  `status` int NOT NULL COMMENT '操作状态',
  `user` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '操作用户',
  `created_time` datetime NOT NULL COMMENT '操作时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '操作记录';

-- ----------------------------
-- Table structure for tangula_store_pool
-- ----------------------------
DROP TABLE IF EXISTS `tangula_store_pool`;
CREATE TABLE `tangula_store_pool` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `uuid` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '唯一标识',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '名称',
  `desc` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '描述',
  `create_user` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '创建用户',
  `created_time` datetime NOT NULL COMMENT '创建时间',
  `updated_time` datetime NOT NULL COMMENT '更新时间',
  `deleted_time` datetime COMMENT '删除时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '存储池信息';

-- ----------------------------
-- Table structure for tangula_image
-- ----------------------------
DROP TABLE IF EXISTS `tangula_image`;
CREATE TABLE `tangula_image` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `uuid` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '唯一标识',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '名称',
  `desc` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '描述',
  `size` bigint DEFAULT 0 NOT NULL COMMENT '镜像大小',
  `type` int DEFAULT 1003 NOT NULL COMMENT '镜像类型',
  `status` int DEFAULT 1024 NOT NULL COMMENT '镜像状态',
  `create_user` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '创建用户',
  `auth_type` int DEFAULT 111 NOT NULL COMMENT '主机权限类型',
  `created_time` datetime NOT NULL COMMENT '创建时间',
  `updated_time` datetime NOT NULL COMMENT '更新时间',
  `deleted_time` datetime COMMENT '删除时间',
  `pool_id` int(11) NOT NULL COMMENT '存储池ID',
  PRIMARY KEY (`id`),
  KEY `pool_id` (`pool_id`),
  CONSTRAINT `tangula_store_pool_ibfk_1` FOREIGN KEY (`pool_id`) REFERENCES `tangula_store_pool` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '镜像信息';

-- ----------------------------
-- Table structure for tangula_replica
-- ----------------------------
DROP TABLE IF EXISTS `tangula_replica`;
CREATE TABLE `tangula_replica` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `uuid` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '唯一标识',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '名称',
  `desc` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '描述',
  `size` bigint DEFAULT 0 NOT NULL COMMENT '副本大小',
  `type` int DEFAULT 1003 NOT NULL COMMENT '副本类型',
  `status` int DEFAULT 1024 NOT NULL COMMENT '副本状态',
  `export` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '导出路径',
  `create_user` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '创建用户',
  `created_time` datetime NOT NULL COMMENT '创建时间',
  `updated_time` datetime NOT NULL COMMENT '更新时间',
  `deleted_time` datetime COMMENT '删除时间',
  `pool_id` int(11) NOT NULL COMMENT '存储池ID',
  PRIMARY KEY (`id`),
  KEY `pool_id` (`pool_id`),
  CONSTRAINT `tangula_store_pool_ibfk_2` FOREIGN KEY (`pool_id`) REFERENCES `tangula_store_pool` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '副本信息';

-- ----------------------------
-- Table structure for tangula_snapshot
-- ----------------------------
DROP TABLE IF EXISTS `tangula_snapshot`;
CREATE TABLE `tangula_snapshot` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `uuid` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '唯一标识',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '名称',
  `desc` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '描述',
  `created_time` datetime NOT NULL COMMENT '创建时间',
  `updated_time` datetime NOT NULL COMMENT '更新时间',
  `deleted_time` datetime COMMENT '删除时间',
  `replica_id` int(11) NOT NULL COMMENT '副本ID',
  PRIMARY KEY (`id`),
  KEY `replica_id` (`replica_id`),
  CONSTRAINT `tangula_replica_ibfk_1` FOREIGN KEY (`replica_id`) REFERENCES `tangula_replica` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '副本快照信息';

-- ----------------------------
-- Table structure for tangula_mount_info
-- ----------------------------
DROP TABLE IF EXISTS `tangula_mount_info`;
CREATE TABLE `tangula_mount_info` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `replica_id` int(11) NOT NULL COMMENT '副本ID',
  `target_type` int DEFAULT 0 NOT NULL COMMENT '应用类型',
  `target_id` int(11) NOT NULL COMMENT '挂载平台ID',
  `mount_param` varchar(8192) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '副本挂载参数',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '副本挂载信息';

-- ----------------------------
-- Table structure for tangula_instance
-- ----------------------------
DROP TABLE IF EXISTS `tangula_instance`;
CREATE TABLE `tangula_instance` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `type` int DEFAULT 1 NOT NULL COMMENT '操作类型',
  `status` int DEFAULT 1 NOT NULL COMMENT '执行状态',
  `replica_id` int(11) NOT NULL COMMENT '副本ID',
  `target_type` int DEFAULT 0 NOT NULL COMMENT '应用类型',
  `target_id` int(11) NOT NULL COMMENT '挂载平台ID',
  `mount_point` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '副本挂载点',
  `created_time` datetime NOT NULL COMMENT '创建时间',
  `updated_time` datetime NOT NULL COMMENT '更新时间',
  `deleted_time` datetime COMMENT '删除时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '实例信息';

-- ----------------------------
-- Table structure for tangula_instance_log
-- ----------------------------
DROP TABLE IF EXISTS `tangula_instance_log`;
CREATE TABLE `tangula_instance_log` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `level` int NOT NULL COMMENT '日志等级',
  `info` varchar(128) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '日志信息',
  `detail` varchar(8192) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '日志详情',
  `instance_id` int(11) NOT NULL COMMENT '操作记录ID',
  `created_time` datetime NOT NULL COMMENT '日志创建时间',
  `updated_time` datetime NOT NULL COMMENT '更新时间',
  `deleted_time` datetime COMMENT '删除时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '实例日志详情';

-- ----------------------------
-- Table structure for tangula_script
-- ----------------------------
DROP TABLE IF EXISTS `tangula_script`;
CREATE TABLE `tangula_script` (
  `id` int(11) NOT NULL AUTO_INCREMENT COMMENT '主键ID',
  `uuid` varchar(32) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '唯一标识',
  `name` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '名称',
  `desc` varchar(512) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '描述',
  `label` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci DEFAULT '' COMMENT '标签',
  `create_user` varchar(50) CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci NOT NULL COMMENT '创建用户',
  `created_time` datetime NOT NULL COMMENT '创建时间',
  `updated_time` datetime NOT NULL COMMENT '更新时间',
  `deleted_time` datetime COMMENT '删除时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT '脚本信息';
