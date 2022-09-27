/*
SQLyog Ultimate v12.09 (64 bit)
MySQL - 5.7.38 : Database - gqbot
*********************************************************************
*/

/*!40101 SET NAMES utf8 */;

/*!40101 SET SQL_MODE=''*/;

/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;
USE `gqbot`;

/*Table structure for table `timed_task` */

DROP TABLE IF EXISTS `timed_task`;

CREATE TABLE `timed_task` (
  `id` int(10) NOT NULL AUTO_INCREMENT,
  `created_id` int(10) NOT NULL COMMENT '创建者id',
  `task_name` varchar(1000) COLLATE utf8mb4_bin NOT NULL COMMENT '定时任务名称',
  `timed_start` int(2) NOT NULL COMMENT '是否定时开始：1-是  2-否',
  `start_time` bigint(20) DEFAULT NULL COMMENT '开始时间',
  `timed_end` int(2) NOT NULL COMMENT '是否定时结束  1-是  2-否',
  `end_time` bigint(20) DEFAULT NULL COMMENT '结束时间',
  `timing_strategy` varchar(1000) COLLATE utf8mb4_bin NOT NULL COMMENT '定时策略',
  `timer_type` int(10) NOT NULL COMMENT '定时器类型',
  `send_to` int(32) DEFAULT NULL COMMENT '发送的群号或者qq号',
  `sent_content` varchar(1000) COLLATE utf8mb4_bin DEFAULT NULL COMMENT '发送的内容',
  `status` int(2) NOT NULL COMMENT '定时器状态 1-待开始 2-运行中 3-已结束',
  `created_time` bigint(20) DEFAULT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=20 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

/*Table structure for table `user` */

DROP TABLE IF EXISTS `user`;

CREATE TABLE `user` (
  `id` int(10) NOT NULL AUTO_INCREMENT COMMENT '用户id',
  `user_name` varchar(64) COLLATE utf8mb4_bin NOT NULL COMMENT '用户名',
  `password` varchar(1000) COLLATE utf8mb4_bin NOT NULL COMMENT '密码',
  `created_time` bigint(20) NOT NULL COMMENT '创建时间',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uniq` (`user_name`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_bin;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;
