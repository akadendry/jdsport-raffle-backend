/*
 Navicat Premium Data Transfer

 Source Server         : Localhost-Docker
 Source Server Type    : MySQL
 Source Server Version : 80027
 Source Host           : localhost:3306
 Source Schema         : raffle

 Target Server Type    : MySQL
 Target Server Version : 80027
 File Encoding         : 65001

 Date: 19/05/2022 13:28:34
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- Table structure for participants
-- ----------------------------
DROP TABLE IF EXISTS `participants`;
CREATE TABLE `participants` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `raffle_id` bigint unsigned DEFAULT NULL,
  `user_id` bigint unsigned DEFAULT NULL,
  `erajaya_club_user_id` bigint unsigned DEFAULT NULL,
  `raffle_product_id` bigint unsigned DEFAULT NULL,
  `raffle_product_size_stock_id` bigint unsigned DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `created_by` varchar(256) DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `updated_by` varchar(256) DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `deleted_by` varchar(256) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Table structure for raffle_product_size_stocks
-- ----------------------------
DROP TABLE IF EXISTS `raffle_product_size_stocks`;
CREATE TABLE `raffle_product_size_stocks` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `raffle_product_id` bigint unsigned DEFAULT NULL,
  `size` varchar(256) DEFAULT NULL,
  `stock` varchar(256) DEFAULT NULL,
  `url_product` varchar(256) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `created_by` varchar(256) DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `updated_by` varchar(256) DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `deleted_by` varchar(256) DEFAULT NULL,
  `sku` varchar(256) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Table structure for raffle_products
-- ----------------------------
DROP TABLE IF EXISTS `raffle_products`;
CREATE TABLE `raffle_products` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `raffle_id` bigint unsigned DEFAULT NULL,
  `name` varchar(256) DEFAULT NULL,
  `description` text,
  `image` varchar(256) DEFAULT NULL,
  `image_mobile` varchar(256) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `created_by` varchar(256) DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `updated_by` varchar(256) DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `deleted_by` varchar(256) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Table structure for raffles
-- ----------------------------
DROP TABLE IF EXISTS `raffles`;
CREATE TABLE `raffles` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `name` varchar(256) DEFAULT NULL,
  `start_date_registration` datetime DEFAULT NULL,
  `end_date_registration` datetime DEFAULT NULL,
  `announcement_date` datetime DEFAULT NULL,
  `start_date_pay` datetime DEFAULT NULL,
  `end_date_pay` datetime DEFAULT NULL,
  `banner` varchar(256) DEFAULT NULL,
  `copyright` varchar(256) DEFAULT NULL,
  `slug` varchar(256) DEFAULT NULL,
  `created_at` datetime DEFAULT NULL,
  `created_by` varchar(256) DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `updated_by` varchar(256) DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `deleted_by` varchar(256) DEFAULT NULL,
  `banner_mobile` varchar(256) DEFAULT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

-- ----------------------------
-- Table structure for users
-- ----------------------------
DROP TABLE IF EXISTS `users`;
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `erajaya_club_user_id` bigint unsigned DEFAULT NULL,
  `first_name` varchar(256) DEFAULT NULL,
  `last_name` varchar(256) DEFAULT NULL,
  `phone` varchar(256) DEFAULT NULL,
  `email` varchar(256) DEFAULT NULL,
  `identity_no` varchar(256) DEFAULT NULL,
  `password` longblob,
  `instagram` varchar(256) DEFAULT NULL,
  `user_type` varchar(256) DEFAULT NULL,
  `role_id` bigint unsigned DEFAULT NULL,
  `status` bigint DEFAULT '1',
  `created_at` datetime DEFAULT NULL,
  `created_by` varchar(256) DEFAULT NULL,
  `updated_at` datetime DEFAULT NULL,
  `updated_by` varchar(256) DEFAULT NULL,
  `deleted_at` datetime DEFAULT NULL,
  `deleted_by` varchar(256) DEFAULT NULL,
  `reason` varchar(256) DEFAULT NULL,
  `token` text,
  PRIMARY KEY (`id`),
  UNIQUE KEY `erajaya_club_user_id` (`erajaya_club_user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

SET FOREIGN_KEY_CHECKS = 1;
