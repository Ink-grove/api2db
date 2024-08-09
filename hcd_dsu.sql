/*
 Navicat Premium Data Transfer

 Source Server         : 127.0.0.1
 Source Server Type    : MySQL
 Source Server Version : 80012
 Source Host           : 127.0.0.1:3309
 Source Schema         : tzzhian_db

 Target Server Type    : MySQL
 Target Server Version : 80012
 File Encoding         : 65001

 Date: 20/09/2023 16:43:12
*/

SET NAMES utf8mb4;
SET
FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for hcd_dsu
-- ----------------------------
DROP TABLE IF EXISTS `hcd_dsu`;
CREATE TABLE `hcd_dsu`
(
    `id`          int(11) NOT NULL AUTO_INCREMENT,
    `project_id`  int(11) DEFAULT NULL COMMENT '项目',
    `sync_type`   int(11) DEFAULT NULL COMMENT '同步器类型',
    `name`        varchar(80) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '任务名称',
    `init_param`  text CHARACTER SET utf8 COLLATE utf8_general_ci COMMENT '初始化参数，主要为数据源、目标数据库以及对应字段映射情况参数',
    `param`       text CHARACTER SET utf8 COLLATE utf8_general_ci COMMENT '任务参数，',
    `status`      tinyint(4) NOT NULL DEFAULT '1' COMMENT '状态',
    `dsu_type`    int(11) DEFAULT NULL COMMENT '数据同步类型：1.【永久运行】  2.【按时执行】依照cron指定的参数  3.【按时执行并且启动的时候执行一次】 4.【仅启动时执行】',
    `dsu_mode`    int(11) DEFAULT NULL COMMENT '数据同步模式：1.【正序运行】  2.【倒序运行】 3.【中序向前运行】  4.【中序向后运行】',
    `spec`        varchar(100) CHARACTER SET utf8 COLLATE utf8_bin                DEFAULT NULL COMMENT '定时任务执行参数',
    `discription` varchar(200) CHARACTER SET utf8 COLLATE utf8_general_ci         DEFAULT NULL COMMENT '任务描述',
    `last_time`   timestamp NULL DEFAULT NULL COMMENT '上次执行时间',
    `create_time` datetime                                               NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` datetime                                               NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`) USING BTREE
) ENGINE=MyISAM AUTO_INCREMENT=23 DEFAULT CHARSET=utf8 COMMENT='数据同步器 - 任务配置表';

-- ----------------------------
-- Records of hcd_dsu
-- ----------------------------
INSERT INTO `tzzhian_db`.`hcd_dsu` (`id`, `project_id`, `sync_type`, `name`, `init_param`, `param`, `status`, `dsu_type`, `dsu_mode`, `spec`, `discription`, `last_time`, `create_time`, `update_time`) VALUES (1, 1, 1, '测试通用api同步', '{\r\n    \"src_config\":{\r\n        \"first_page_num\":1,\r\n        \"page_size\":100,\r\n        \"page_count\":0\r\n    },\r\n    \"target_config\":{\r\n        \"table_name\":\"tz_private_db.test_build_data\",\r\n        \"page_size\":0\r\n    },\r\n    \"filed\":{\r\n        \"unique_field_name\":\"id\",\r\n        \"md5_filed_list\":\"mc\",\r\n        \"save_filed_list\":\"\",\r\n        \"filed_alias\":{\r\n            \"mc\":\"name\",\r\n            \"x\":\"lng\",\r\n            \"y\":\"lat\"\r\n        }\r\n    }\r\n}', '{\r\n    \"http_method\":\"post\",\r\n    \"call_method\":\"ApiReq_Normal\",\r\n    \"path\":\"/api/tz/getBuildList\",\r\n    \"body\":{\r\n        \"parameterSetObject\":{}\r\n    },\r\n    \"data_filed\":\"data\",\r\n    \"page_count_filed\":\"allpage\",\r\n    \"page_filed_position\":{\r\n       \"type\":1,\r\n       \"num_position\":\"pageNo\",\r\n       \"count_position\":\"pageSize\"\r\n   },\r\n    \"request_info\":{\r\n       \"base_server\":\"http://127.0.0.1:12144\"\r\n   }\r\n}', 1, 4, NULL, '0 0,30 7,8,9,10,11,12,13,14,15,16,17,18 * * ?', NULL, NULL, '2023-09-10 20:28:37', '2023-09-22 14:57:18');
INSERT INTO `tzzhian_db`.`hcd_dsu` (`id`, `project_id`, `sync_type`, `name`, `init_param`, `param`, `status`, `dsu_type`, `dsu_mode`, `spec`, `discription`, `last_time`, `create_time`, `update_time`) VALUES (2, 1, 1, '测试通用api同步', '{\r\n    \"src_config\":{\r\n        \"first_page_num\":1,\r\n        \"page_size\":100,\r\n        \"page_count\":0\r\n    },\r\n    \"target_config\":{\r\n        \"table_name\":\"tz_private_db.test_build_data_1\",\r\n        \"page_size\":0\r\n    },\r\n    \"filed\":{\r\n        \"unique_field_name\":\"id\",\r\n        \"md5_filed_list\":\"mc\",\r\n        \"save_filed_list\":\"\",\r\n        \"filed_alias\":{\r\n            \"mc\":\"name\",\r\n            \"x\":\"lng\",\r\n            \"y\":\"lat\"\r\n        }\r\n    }\r\n}', '{\r\n    \"http_method\":\"post\",\r\n    \"call_method\":\"ApiReq_Normal\",\r\n    \"path\":\"/api/tz/getBuildList\",\r\n    \"body\":{\r\n        \"parameterSetObject\":{}\r\n    },\r\n    \"data_filed\":\"data\",\r\n    \"page_count_filed\":\"allpage\",\r\n    \"page_filed_position\":{\r\n       \"type\":1,\r\n       \"num_position\":\"pageNo\",\r\n       \"count_position\":\"pageSize\"\r\n   },\r\n    \"request_info\":{\r\n       \"base_server\":\"http://127.0.0.1:12144\"\r\n   }\r\n}', 0, 4, NULL, NULL, NULL, NULL, '2023-09-10 20:28:37', '2023-09-22 15:47:18');


-- ----------------------------
-- Table structure for hcd_dsu_fail_task
-- ----------------------------
DROP TABLE IF EXISTS `hcd_dsu_fail_task`;
CREATE TABLE `hcd_dsu_fail_task`
(
    `id`          int(11) NOT NULL AUTO_INCREMENT,
    `project_id`  int(11) NULL DEFAULT NULL COMMENT '项目',
    `sync_type`   int(11) NULL DEFAULT NULL COMMENT '同步器类型',
    `name`        varchar(80) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '任务名称',
    `task_id`     int(11) NULL DEFAULT NULL COMMENT '任务id',
    `call_info`   text CHARACTER SET utf8 COLLATE utf8_general_ci NULL COMMENT '调用信息，保存请求的path，body以及服务器参数等，用于任务补偿',
    `discription` varchar(200) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '任务描述',
    `create_time` datetime                                               NOT NULL DEFAULT CURRENT_TIMESTAMP,
    `update_time` datetime                                               NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (`id`) USING BTREE
) ENGINE = MyISAM AUTO_INCREMENT = 3 CHARACTER SET = utf8 COLLATE = utf8_general_ci COMMENT = '数据同步器 - 执行任务失败表' ROW_FORMAT = Dynamic;

-- ----------------------------
-- Records of hcd_dsu_fail_task
-- ----------------------------
INSERT INTO `hcd_dsu_fail_task`
VALUES (1, 1, 1, '测试通用api同步', 1,
        '{\"client_param\":\"{\r\n\"base_server\":\"https://127.0.0.1:12144\"\r\n}\",\"input_param\":{\"ApiParam\":{\"Body\":{\"pageNo\":4,\"pageSize\":100,\"parameterSetObject\":{}},\"HttpMethod\":\"POST\",\"Path\":\"/api/tz/getBuildList\"},\"DbParam\":null,\"PageSize\":100}}',
        'Post \"https://127.0.0.1:12144/api/tz/getBuildList\": dial tcp 127.0.0.1:12144: i/o timeout',
        '2023-09-14 17:10:30', '2023-09-14 17:10:30');
INSERT INTO `hcd_dsu_fail_task`
VALUES (2, 1, 1, '测试通用api同步', 1,
        '{\"client_param\":\"{\r\n\"base_server\":\"https://127.0.0.1:12144\"\r\n}\",\"input_param\":{\"ApiParam\":{\"Body\":{\"pageNo\":6,\"pageSize\":100,\"parameterSetObject\":{}},\"HttpMethod\":\"POST\",\"Path\":\"/api/tz/getBuildList\"},\"DbParam\":null,\"PageSize\":100}}',
        'Post \"https://127.0.0.1:12144/api/tz/getBuildList\": dial tcp 127.0.0.1:12144: i/o timeout',
        '2023-09-14 17:11:17', '2023-09-14 17:11:17');

SET
FOREIGN_KEY_CHECKS = 1;
