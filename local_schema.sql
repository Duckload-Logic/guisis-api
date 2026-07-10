-- MySQL dump 10.13  Distrib 8.0.43, for Linux (x86_64)
--
-- Host: localhost    Database: guidance_db
-- ------------------------------------------------------
-- Server version	8.0.43

/*!40101 SET @OLD_CHARACTER_SET_CLIENT=@@CHARACTER_SET_CLIENT */;
/*!40101 SET @OLD_CHARACTER_SET_RESULTS=@@CHARACTER_SET_RESULTS */;
/*!40101 SET @OLD_COLLATION_CONNECTION=@@COLLATION_CONNECTION */;
/*!50503 SET NAMES utf8mb4 */;
/*!40103 SET @OLD_TIME_ZONE=@@TIME_ZONE */;
/*!40103 SET TIME_ZONE='+00:00' */;
/*!40014 SET @OLD_UNIQUE_CHECKS=@@UNIQUE_CHECKS, UNIQUE_CHECKS=0 */;
/*!40014 SET @OLD_FOREIGN_KEY_CHECKS=@@FOREIGN_KEY_CHECKS, FOREIGN_KEY_CHECKS=0 */;
/*!40101 SET @OLD_SQL_MODE=@@SQL_MODE, SQL_MODE='NO_AUTO_VALUE_ON_ZERO' */;
/*!40111 SET @OLD_SQL_NOTES=@@SQL_NOTES, SQL_NOTES=0 */;

--
-- Table structure for table `academic_settings`
--

DROP TABLE IF EXISTS `academic_settings`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `academic_settings` (
  `id` int NOT NULL DEFAULT '1',
  `current_year_start` int NOT NULL,
  `current_year_end` int NOT NULL,
  `current_term` int NOT NULL,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  CONSTRAINT `chk_academic_settings_id` CHECK ((`id` = 1))
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `activity_options`
--

DROP TABLE IF EXISTS `activity_options`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `activity_options` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `category` enum('academic','extra_curricular','both') NOT NULL,
  `is_active` tinyint(1) DEFAULT '1',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=15 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `addresses`
--

DROP TABLE IF EXISTS `addresses`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `addresses` (
  `id` int NOT NULL AUTO_INCREMENT,
  `region_code` varchar(10) DEFAULT NULL,
  `province_code` varchar(10) DEFAULT NULL,
  `city_code` varchar(10) DEFAULT NULL,
  `barangay_code` varchar(10) NOT NULL,
  `street_detail` varchar(255) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_addresses_region_id` (`region_code`),
  KEY `idx_addresses_province_id` (`province_code`),
  KEY `idx_addresses_city_id` (`city_code`),
  KEY `idx_addresses_barangay_id` (`barangay_code`),
  CONSTRAINT `addresses_ibfk_1` FOREIGN KEY (`region_code`) REFERENCES `regions` (`code`),
  CONSTRAINT `addresses_ibfk_2` FOREIGN KEY (`city_code`) REFERENCES `cities` (`code`),
  CONSTRAINT `addresses_ibfk_3` FOREIGN KEY (`province_code`) REFERENCES `provinces` (`code`),
  CONSTRAINT `addresses_ibfk_4` FOREIGN KEY (`barangay_code`) REFERENCES `barangays` (`code`)
) ENGINE=InnoDB AUTO_INCREMENT=101 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `admission_slip_categories`
--

DROP TABLE IF EXISTS `admission_slip_categories`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `admission_slip_categories` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`),
  UNIQUE KEY `unique_idx_admission_slip_category_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `admission_slips`
--

DROP TABLE IF EXISTS `admission_slips`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `admission_slips` (
  `id` char(36) NOT NULL,
  `iir_id` char(36) NOT NULL,
  `category_id` int NOT NULL,
  `reason` text NOT NULL,
  `date_of_absence` date NOT NULL,
  `date_needed` date NOT NULL,
  `status_id` int NOT NULL DEFAULT '1',
  `admin_notes` text,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_admission_slips_status_id` (`status_id`),
  KEY `idx_admission_slips_category_id` (`category_id`),
  KEY `idx_admission_slips_iir_id` (`iir_id`),
  CONSTRAINT `admission_slips_ibfk_2` FOREIGN KEY (`status_id`) REFERENCES `statuses` (`id`),
  CONSTRAINT `admission_slips_ibfk_3` FOREIGN KEY (`category_id`) REFERENCES `admission_slip_categories` (`id`),
  CONSTRAINT `fk_admission_slips_iir` FOREIGN KEY (`iir_id`) REFERENCES `iir_records` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `admission_tickets`
--

DROP TABLE IF EXISTS `admission_tickets`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `admission_tickets` (
  `id` char(36) NOT NULL,
  `admission_slip_id` char(36) NOT NULL,
  `ticket_code` varchar(20) NOT NULL,
  `is_verified` tinyint(1) DEFAULT '0',
  `verified_at` timestamp NULL DEFAULT NULL,
  `verified_by` char(36) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `ticket_code` (`ticket_code`),
  KEY `fk_tickets_verified_by` (`verified_by`),
  KEY `idx_tickets_slip_id` (`admission_slip_id`),
  KEY `idx_tickets_code` (`ticket_code`),
  CONSTRAINT `fk_tickets_slip` FOREIGN KEY (`admission_slip_id`) REFERENCES `admission_slips` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_tickets_verified_by` FOREIGN KEY (`verified_by`) REFERENCES `users` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `appointment_categories`
--

DROP TABLE IF EXISTS `appointment_categories`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `appointment_categories` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`),
  UNIQUE KEY `unique_idx_appointment_category_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `appointments`
--

DROP TABLE IF EXISTS `appointments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `appointments` (
  `id` char(36) NOT NULL,
  `iir_id` char(36) DEFAULT NULL,
  `time_slot_id` int NOT NULL,
  `when_date` date NOT NULL,
  `reason` text,
  `admin_notes` text,
  `appointment_category_id` int NOT NULL,
  `urgency_level` enum('LOW','MEDIUM','HIGH','CRITICAL') NOT NULL DEFAULT 'MEDIUM',
  `urgency_score` float NOT NULL,
  `status_id` int NOT NULL DEFAULT '1',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `preferred_date_1` date DEFAULT NULL,
  `preferred_time_slot_id_1` int DEFAULT NULL,
  `preferred_date_2` date DEFAULT NULL,
  `preferred_time_slot_id_2` int DEFAULT NULL,
  `preferred_date_3` date DEFAULT NULL,
  `preferred_time_slot_id_3` int DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_appointments_time_slot_id` (`time_slot_id`),
  KEY `idx_appointments_appointment_category_id` (`appointment_category_id`),
  KEY `idx_appointments_status_id` (`status_id`),
  KEY `idx_appointments_when_date` (`when_date`),
  KEY `idx_appointments_iir_id` (`iir_id`),
  KEY `fk_appointments_pref_time_slot_1` (`preferred_time_slot_id_1`),
  KEY `fk_appointments_pref_time_slot_2` (`preferred_time_slot_id_2`),
  KEY `fk_appointments_pref_time_slot_3` (`preferred_time_slot_id_3`),
  CONSTRAINT `appointments_ibfk_2` FOREIGN KEY (`time_slot_id`) REFERENCES `time_slots` (`id`),
  CONSTRAINT `appointments_ibfk_3` FOREIGN KEY (`status_id`) REFERENCES `statuses` (`id`),
  CONSTRAINT `appointments_ibfk_4` FOREIGN KEY (`appointment_category_id`) REFERENCES `appointment_categories` (`id`),
  CONSTRAINT `fk_appointments_iir` FOREIGN KEY (`iir_id`) REFERENCES `iir_records` (`id`) ON DELETE SET NULL,
  CONSTRAINT `fk_appointments_pref_time_slot_1` FOREIGN KEY (`preferred_time_slot_id_1`) REFERENCES `time_slots` (`id`) ON DELETE SET NULL,
  CONSTRAINT `fk_appointments_pref_time_slot_2` FOREIGN KEY (`preferred_time_slot_id_2`) REFERENCES `time_slots` (`id`) ON DELETE SET NULL,
  CONSTRAINT `fk_appointments_pref_time_slot_3` FOREIGN KEY (`preferred_time_slot_id_3`) REFERENCES `time_slots` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `barangays`
--

DROP TABLE IF EXISTS `barangays`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `barangays` (
  `id` int NOT NULL AUTO_INCREMENT,
  `code` varchar(10) NOT NULL,
  `name` varchar(100) DEFAULT NULL,
  `city_code` varchar(10) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `code` (`code`),
  UNIQUE KEY `unique_idx_barangay_code` (`code`),
  UNIQUE KEY `unique_idx_city_barangay` (`city_code`,`name`),
  KEY `idx_barangays_city_code` (`city_code`),
  CONSTRAINT `barangays_ibfk_1` FOREIGN KEY (`city_code`) REFERENCES `cities` (`code`)
) ENGINE=InnoDB AUTO_INCREMENT=42028 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `cities`
--

DROP TABLE IF EXISTS `cities`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `cities` (
  `id` int NOT NULL AUTO_INCREMENT,
  `code` varchar(10) NOT NULL,
  `name` varchar(100) DEFAULT NULL,
  `type` varchar(20) DEFAULT NULL,
  `zip_code` varchar(10) DEFAULT NULL,
  `district` varchar(50) DEFAULT NULL,
  `province_code` varchar(10) DEFAULT NULL,
  `region_code` varchar(10) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `code` (`code`),
  UNIQUE KEY `unique_idx_city_code` (`code`),
  UNIQUE KEY `unique_idx_province_city` (`province_code`,`name`),
  KEY `idx_cities_region_code` (`region_code`),
  KEY `idx_cities_province_code` (`province_code`),
  CONSTRAINT `cities_ibfk_1` FOREIGN KEY (`province_code`) REFERENCES `provinces` (`code`),
  CONSTRAINT `cities_ibfk_2` FOREIGN KEY (`region_code`) REFERENCES `regions` (`code`)
) ENGINE=InnoDB AUTO_INCREMENT=1657 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `civil_status_types`
--

DROP TABLE IF EXISTS `civil_status_types`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `civil_status_types` (
  `id` int NOT NULL,
  `status_name` varchar(50) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `status_name` (`status_name`),
  UNIQUE KEY `unique_idx_civil_status_name` (`status_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `courses`
--

DROP TABLE IF EXISTS `courses`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `courses` (
  `id` int NOT NULL AUTO_INCREMENT,
  `code` varchar(20) NOT NULL,
  `course_name` varchar(100) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `code` (`code`),
  UNIQUE KEY `course_name` (`course_name`),
  UNIQUE KEY `unique_idx_course_code` (`code`),
  UNIQUE KEY `unique_idx_course_name` (`course_name`)
) ENGINE=InnoDB AUTO_INCREMENT=12 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `educational_attainments`
--

DROP TABLE IF EXISTS `educational_attainments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `educational_attainments` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(255) NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `educational_backgrounds`
--

DROP TABLE IF EXISTS `educational_backgrounds`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `educational_backgrounds` (
  `id` int NOT NULL AUTO_INCREMENT,
  `iir_id` char(36) NOT NULL,
  `nature_of_schooling` enum('Continuous','Interrupted') NOT NULL,
  `interrupted_details` varchar(255) DEFAULT NULL,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_educational_backgrounds_iir_id` (`iir_id`),
  CONSTRAINT `educational_backgrounds_ibfk_1` FOREIGN KEY (`iir_id`) REFERENCES `iir_records` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=51 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `educational_levels`
--

DROP TABLE IF EXISTS `educational_levels`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `educational_levels` (
  `id` int NOT NULL AUTO_INCREMENT,
  `level_name` varchar(50) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `level_name` (`level_name`),
  UNIQUE KEY `unique_idx_level_name` (`level_name`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `emergency_contacts`
--

DROP TABLE IF EXISTS `emergency_contacts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `emergency_contacts` (
  `id` int NOT NULL AUTO_INCREMENT,
  `iir_id` char(36) NOT NULL,
  `first_name` varchar(100) NOT NULL,
  `middle_name` varchar(100) DEFAULT NULL,
  `last_name` varchar(100) NOT NULL,
  `suffix_name` varchar(50) DEFAULT NULL,
  `contact_number` varchar(20) NOT NULL,
  `relationship_id` int NOT NULL,
  `address_id` int NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_emergency_contacts_iir_id` (`iir_id`),
  KEY `idx_emergency_contacts_relationship_id` (`relationship_id`),
  KEY `idx_emergency_contacts_address_id` (`address_id`),
  CONSTRAINT `emergency_contacts_ibfk_1` FOREIGN KEY (`iir_id`) REFERENCES `iir_records` (`id`) ON DELETE CASCADE,
  CONSTRAINT `emergency_contacts_ibfk_2` FOREIGN KEY (`relationship_id`) REFERENCES `student_relationship_types` (`id`),
  CONSTRAINT `emergency_contacts_ibfk_3` FOREIGN KEY (`address_id`) REFERENCES `addresses` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=51 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `family_backgrounds`
--

DROP TABLE IF EXISTS `family_backgrounds`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `family_backgrounds` (
  `id` int NOT NULL AUTO_INCREMENT,
  `iir_id` char(36) NOT NULL,
  `parental_status_id` int NOT NULL,
  `parental_status_details` varchar(255) DEFAULT NULL,
  `brothers` int NOT NULL,
  `sisters` int NOT NULL,
  `employed_siblings` int NOT NULL,
  `ordinal_position` int NOT NULL,
  `have_quiet_place_to_study` tinyint(1) NOT NULL,
  `is_sharing_room` tinyint(1) NOT NULL,
  `room_sharing_details` varchar(255) DEFAULT NULL,
  `nature_of_residence_id` int NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `iir_id` (`iir_id`),
  KEY `idx_family_backgrounds_iir_id` (`iir_id`),
  KEY `idx_family_backgrounds_parental_status_id` (`parental_status_id`),
  KEY `idx_family_backgrounds_nature_of_residence_id` (`nature_of_residence_id`),
  CONSTRAINT `family_backgrounds_ibfk_1` FOREIGN KEY (`iir_id`) REFERENCES `iir_records` (`id`) ON DELETE CASCADE,
  CONSTRAINT `family_backgrounds_ibfk_2` FOREIGN KEY (`nature_of_residence_id`) REFERENCES `nature_of_residence_types` (`id`),
  CONSTRAINT `family_backgrounds_ibfk_3` FOREIGN KEY (`parental_status_id`) REFERENCES `parental_status_types` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=51 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `file_ocr_results`
--

DROP TABLE IF EXISTS `file_ocr_results`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `file_ocr_results` (
  `file_id` char(36) NOT NULL,
  `raw_text` longtext,
  `structured_data` json DEFAULT NULL,
  `engine_v` varchar(50) DEFAULT NULL,
  `confidence` float DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`file_id`),
  CONSTRAINT `file_ocr_results_ibfk_1` FOREIGN KEY (`file_id`) REFERENCES `files` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `files`
--

DROP TABLE IF EXISTS `files`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `files` (
  `id` char(36) NOT NULL,
  `file_name` varchar(255) NOT NULL,
  `file_url` varchar(255) NOT NULL,
  `file_type` varchar(255) NOT NULL,
  `file_size` bigint NOT NULL,
  `mime_type` varchar(255) NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `deleted_at` timestamp NULL DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `idx_files_id` (`id`),
  KEY `idx_files_created_at` (`created_at`),
  KEY `idx_files_deleted_at` (`deleted_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;



--
-- Table structure for table `iir_drafts`
--

DROP TABLE IF EXISTS `iir_drafts`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `iir_drafts` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` char(36) NOT NULL,
  `data` json NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_iir_drafts_user_unique` (`user_id`),
  KEY `idx_iir_drafts_user_id` (`user_id`),
  CONSTRAINT `fk_iir_drafts_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=59 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `iir_records`
--

DROP TABLE IF EXISTS `iir_records`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `iir_records` (
  `id` char(36) NOT NULL,
  `user_id` char(36) NOT NULL,
  `is_submitted` tinyint(1) DEFAULT '0',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_iir_records_user_unique` (`user_id`),
  KEY `idx_iir_records_user_id` (`user_id`),
  CONSTRAINT `fk_iir_records_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `income_ranges`
--

DROP TABLE IF EXISTS `income_ranges`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `income_ranges` (
  `id` int NOT NULL AUTO_INCREMENT,
  `range_text` varchar(100) NOT NULL,
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=11 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `m2m_clients`
--

DROP TABLE IF EXISTS `m2m_clients`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `m2m_clients` (
  `id` int NOT NULL AUTO_INCREMENT,
  `user_id` char(36) NOT NULL,
  `client_name` varchar(100) NOT NULL,
  `client_id` varchar(36) NOT NULL,
  `client_secret_hash` varchar(64) NOT NULL,
  `client_description` varchar(255) NOT NULL,
  `scopes` json DEFAULT NULL,
  `is_active` tinyint(1) NOT NULL DEFAULT '1',
  `is_verified` tinyint(1) NOT NULL DEFAULT '0',
  `last_used_at` timestamp NULL DEFAULT NULL,
  `expires_at` timestamp NULL DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `has_personal_info_access` tinyint(1) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`),
  UNIQUE KEY `client_id` (`client_id`),
  UNIQUE KEY `unique_idx_m2m_clients_client_id` (`client_id`),
  KEY `idx_m2m_clients_user_id` (`user_id`),
  KEY `idx_m2m_clients_is_active` (`is_active`),
  CONSTRAINT `fk_m2m_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `nature_of_residence_types`
--

DROP TABLE IF EXISTS `nature_of_residence_types`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `nature_of_residence_types` (
  `id` int NOT NULL AUTO_INCREMENT,
  `residence_type_name` varchar(100) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `residence_type_name` (`residence_type_name`),
  UNIQUE KEY `unique_idx_residence_type_name` (`residence_type_name`)
) ENGINE=InnoDB AUTO_INCREMENT=8 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `notifications`
--

DROP TABLE IF EXISTS `notifications`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `notifications` (
  `id` char(36) NOT NULL,
  `receiver_id` char(36) NOT NULL,
  `actor_id` char(36) DEFAULT NULL,
  `target_id` char(36) DEFAULT NULL,
  `target_type` varchar(50) DEFAULT NULL,
  `title` varchar(255) NOT NULL,
  `message` text NOT NULL,
  `type` enum('Appointment','Slip','Guidance','System','General') DEFAULT 'System',
  `is_read` tinyint(1) DEFAULT '0',
  `is_touched` tinyint(1) DEFAULT '0',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_notifications_receiver_id` (`receiver_id`),
  KEY `idx_notifications_actor_id` (`actor_id`),
  CONSTRAINT `fk_notifications_actor` FOREIGN KEY (`actor_id`) REFERENCES `users` (`id`) ON DELETE SET NULL,
  CONSTRAINT `fk_notifications_receiver` FOREIGN KEY (`receiver_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `parental_status_types`
--

DROP TABLE IF EXISTS `parental_status_types`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `parental_status_types` (
  `id` int NOT NULL,
  `status_name` varchar(50) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `status_name` (`status_name`),
  UNIQUE KEY `unique_idx_parental_status_name` (`status_name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `profile_pictures`
--

DROP TABLE IF EXISTS `profile_pictures`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `profile_pictures` (
  `file_id` char(36) NOT NULL,
  `user_id` char(36) NOT NULL,
  PRIMARY KEY (`file_id`),
  UNIQUE KEY `user_id` (`user_id`),
  KEY `idx_profile_pictures_user_id` (`user_id`),
  CONSTRAINT `profile_pictures_ibfk_1` FOREIGN KEY (`file_id`) REFERENCES `files` (`id`) ON DELETE CASCADE,
  CONSTRAINT `profile_pictures_ibfk_2` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `provinces`
--

DROP TABLE IF EXISTS `provinces`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `provinces` (
  `id` int NOT NULL AUTO_INCREMENT,
  `code` varchar(10) NOT NULL,
  `name` varchar(100) NOT NULL,
  `region_code` varchar(10) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `code` (`code`),
  UNIQUE KEY `unique_idx_progince_code` (`code`),
  KEY `idx_province_region_code` (`region_code`),
  CONSTRAINT `provinces_ibfk_1` FOREIGN KEY (`region_code`) REFERENCES `regions` (`code`)
) ENGINE=InnoDB AUTO_INCREMENT=83 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `push_subscriptions`
--

DROP TABLE IF EXISTS `push_subscriptions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `push_subscriptions` (
  `id` varchar(36) NOT NULL,
  `user_id` varchar(36) NOT NULL,
  `endpoint` varchar(512) NOT NULL,
  `p256dh_key` text NOT NULL,
  `auth_key` text NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `endpoint` (`endpoint`),
  KEY `idx_push_sub_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `regions`
--

DROP TABLE IF EXISTS `regions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `regions` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(100) NOT NULL,
  `code` varchar(10) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`),
  UNIQUE KEY `code` (`code`),
  UNIQUE KEY `unique_idx_region_name` (`name`),
  UNIQUE KEY `unique_idx_region_code` (`code`)
) ENGINE=InnoDB AUTO_INCREMENT=18 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `related_persons`
--

DROP TABLE IF EXISTS `related_persons`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `related_persons` (
  `id` int NOT NULL AUTO_INCREMENT,
  `first_name` varchar(100) NOT NULL,
  `middle_name` varchar(100) DEFAULT NULL,
  `last_name` varchar(100) NOT NULL,
  `suffix_name` varchar(50) DEFAULT NULL,
  `date_of_birth` date NOT NULL,
  `occupation` varchar(100) DEFAULT NULL,
  `employer_name` varchar(150) DEFAULT NULL,
  `employer_address` varchar(255) DEFAULT NULL,
  `contact_number` varchar(20) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `educational_attainment_id` int DEFAULT NULL,
  PRIMARY KEY (`id`),
  KEY `fk_related_persons_attainment` (`educational_attainment_id`),
  CONSTRAINT `fk_related_persons_attainment` FOREIGN KEY (`educational_attainment_id`) REFERENCES `educational_attainments` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=118 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `religions`
--

DROP TABLE IF EXISTS `religions`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `religions` (
  `id` int NOT NULL AUTO_INCREMENT,
  `religion_name` varchar(100) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `religion_name` (`religion_name`),
  UNIQUE KEY `unique_idx_religion_name` (`religion_name`)
) ENGINE=InnoDB AUTO_INCREMENT=17 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `roles`
--

DROP TABLE IF EXISTS `roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `roles` (
  `id` int NOT NULL,
  `name` varchar(50) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `schema_migrations`
--

DROP TABLE IF EXISTS `schema_migrations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `schema_migrations` (
  `version` bigint NOT NULL,
  `dirty` tinyint(1) NOT NULL,
  PRIMARY KEY (`version`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `school_details`
--

DROP TABLE IF EXISTS `school_details`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `school_details` (
  `id` int NOT NULL AUTO_INCREMENT,
  `eb_id` int NOT NULL,
  `educational_level_id` int NOT NULL,
  `school_name` varchar(255) NOT NULL,
  `school_address` varchar(255) NOT NULL,
  `school_type` enum('Private','Public') NOT NULL,
  `year_started` smallint NOT NULL,
  `year_completed` smallint NOT NULL,
  `awards` text,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_school_details_eb_id` (`eb_id`),
  KEY `idx_school_details_educational_level_id` (`educational_level_id`),
  CONSTRAINT `school_details_ibfk_1` FOREIGN KEY (`eb_id`) REFERENCES `educational_backgrounds` (`id`) ON DELETE CASCADE,
  CONSTRAINT `school_details_ibfk_2` FOREIGN KEY (`educational_level_id`) REFERENCES `educational_levels` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=151 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `seed_migrations`
--

DROP TABLE IF EXISTS `seed_migrations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `seed_migrations` (
  `version` bigint NOT NULL,
  `dirty` tinyint(1) NOT NULL,
  PRIMARY KEY (`version`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `sibling_support_types`
--

DROP TABLE IF EXISTS `sibling_support_types`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `sibling_support_types` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=4 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `significant_notes`
--

DROP TABLE IF EXISTS `significant_notes`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `significant_notes` (
  `id` char(36) NOT NULL,
  `iir_id` char(36) DEFAULT NULL,
  `appointment_id` char(36) DEFAULT NULL,
  `admission_slip_id` char(36) DEFAULT NULL,
  `note` text,
  `remarks` text,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_significant_notes_iir_id` (`iir_id`),
  KEY `idx_significant_notes_appointment_id` (`appointment_id`),
  KEY `idx_significant_notes_admission_slip_id` (`admission_slip_id`),
  CONSTRAINT `fk_sig_notes_admission_slip` FOREIGN KEY (`admission_slip_id`) REFERENCES `admission_slips` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_sig_notes_appointment` FOREIGN KEY (`appointment_id`) REFERENCES `appointments` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_sig_notes_iir` FOREIGN KEY (`iir_id`) REFERENCES `iir_records` (`id`) ON DELETE SET NULL
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `slip_attachments`
--

DROP TABLE IF EXISTS `slip_attachments`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `slip_attachments` (
  `file_id` char(36) NOT NULL,
  `admission_slip_id` char(36) NOT NULL,
  `attachment_type` enum('MEDICAL','EXCUSE LETTER','PARENT VALID ID','OTHER') NOT NULL,
  PRIMARY KEY (`file_id`),
  KEY `idx_slip_attachments_admission_slip_id` (`admission_slip_id`),
  CONSTRAINT `slip_attachments_ibfk_1` FOREIGN KEY (`file_id`) REFERENCES `files` (`id`) ON DELETE CASCADE,
  CONSTRAINT `slip_attachments_ibfk_2` FOREIGN KEY (`admission_slip_id`) REFERENCES `admission_slips` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `statuses`
--

DROP TABLE IF EXISTS `statuses`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `statuses` (
  `id` int NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL,
  `status_type` enum('appointment','slip','both') NOT NULL DEFAULT 'both',
  PRIMARY KEY (`id`),
  UNIQUE KEY `name` (`name`),
  UNIQUE KEY `unique_idx_status_name` (`name`)
) ENGINE=InnoDB AUTO_INCREMENT=10 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `student_activities`
--

DROP TABLE IF EXISTS `student_activities`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `student_activities` (
  `id` int NOT NULL AUTO_INCREMENT,
  `iir_id` char(36) NOT NULL,
  `option_id` int DEFAULT NULL,
  `other_specification` varchar(255) DEFAULT NULL,
  `role` enum('Officer','Member','Other') DEFAULT 'Member',
  `role_specification` varchar(255) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_student_activities_iir_id` (`iir_id`),
  KEY `idx_student_activities_option_id` (`option_id`),
  CONSTRAINT `student_activities_ibfk_1` FOREIGN KEY (`iir_id`) REFERENCES `iir_records` (`id`) ON DELETE CASCADE,
  CONSTRAINT `student_activities_ibfk_2` FOREIGN KEY (`option_id`) REFERENCES `activity_options` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=79 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `student_addresses`
--

DROP TABLE IF EXISTS `student_addresses`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `student_addresses` (
  `id` int NOT NULL AUTO_INCREMENT,
  `iir_id` char(36) NOT NULL,
  `address_id` int NOT NULL,
  `address_type` enum('Residential','Provincial') NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_student_addresses_iir_id` (`iir_id`),
  KEY `idx_student_addresses_address_id` (`address_id`),
  CONSTRAINT `student_addresses_ibfk_1` FOREIGN KEY (`iir_id`) REFERENCES `iir_records` (`id`) ON DELETE CASCADE,
  CONSTRAINT `student_addresses_ibfk_2` FOREIGN KEY (`address_id`) REFERENCES `addresses` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=101 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `student_consultations`
--

DROP TABLE IF EXISTS `student_consultations`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `student_consultations` (
  `id` int NOT NULL AUTO_INCREMENT,
  `iir_id` char(36) NOT NULL,
  `professional_type` enum('Psychiatrist','Psychologist','Counselor') NOT NULL,
  `has_consulted` tinyint(1) DEFAULT '0',
  `when_date` varchar(100) DEFAULT NULL,
  `for_what` text,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_student_consultations_iir_id` (`iir_id`),
  CONSTRAINT `student_consultations_ibfk_1` FOREIGN KEY (`iir_id`) REFERENCES `iir_records` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=48 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `student_cors`
--

DROP TABLE IF EXISTS `student_cors`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `student_cors` (
  `file_id` char(36) NOT NULL,
  `student_id` char(36) NOT NULL,
  `student_number` varchar(20) NOT NULL,
  `course_desc` varchar(255) NOT NULL,
  `course_code` varchar(10) NOT NULL,
  `year_level` int NOT NULL,
  `section` int NOT NULL,
  `campus` varchar(20) NOT NULL,
  `year_start` int NOT NULL DEFAULT '2025',
  `year_end` int NOT NULL DEFAULT '2026',
  `term` int NOT NULL,
  `valid_from` date DEFAULT NULL,
  `valid_until` date DEFAULT NULL,
  PRIMARY KEY (`file_id`),
  KEY `idx_student_cors_student_id` (`student_id`),
  CONSTRAINT `student_cors_ibfk_1` FOREIGN KEY (`file_id`) REFERENCES `files` (`id`) ON DELETE CASCADE,
  CONSTRAINT `student_cors_ibfk_2` FOREIGN KEY (`student_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `student_finances`
--

DROP TABLE IF EXISTS `student_finances`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `student_finances` (
  `id` int NOT NULL AUTO_INCREMENT,
  `iir_id` char(36) NOT NULL,
  `monthly_family_income_range_id` int DEFAULT NULL,
  `other_income_details` varchar(50) DEFAULT NULL,
  `weekly_allowance` decimal(10,2) DEFAULT '0.00',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_idx_student_finances_iir_id` (`iir_id`),
  KEY `idx_student_finances_iir_id` (`iir_id`),
  KEY `idx_student_finances_income_range_id` (`monthly_family_income_range_id`),
  CONSTRAINT `student_finances_ibfk_1` FOREIGN KEY (`iir_id`) REFERENCES `iir_records` (`id`) ON DELETE CASCADE,
  CONSTRAINT `student_finances_ibfk_2` FOREIGN KEY (`monthly_family_income_range_id`) REFERENCES `income_ranges` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=51 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `student_financial_supports`
--

DROP TABLE IF EXISTS `student_financial_supports`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `student_financial_supports` (
  `sf_id` int NOT NULL,
  `support_type_id` int NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `idx_student_financial_supports_sf_id` (`sf_id`),
  KEY `idx_student_financial_supports_support_type_id` (`support_type_id`),
  CONSTRAINT `student_financial_supports_ibfk_1` FOREIGN KEY (`support_type_id`) REFERENCES `student_support_types` (`id`),
  CONSTRAINT `student_financial_supports_ibfk_2` FOREIGN KEY (`sf_id`) REFERENCES `student_finances` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `student_health_records`
--

DROP TABLE IF EXISTS `student_health_records`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `student_health_records` (
  `id` int NOT NULL AUTO_INCREMENT,
  `iir_id` char(36) NOT NULL,
  `vision_has_problem` tinyint(1) DEFAULT '0',
  `vision_details` varchar(255) DEFAULT NULL,
  `hearing_has_problem` tinyint(1) DEFAULT '0',
  `hearing_details` varchar(255) DEFAULT NULL,
  `speech_has_problem` tinyint(1) DEFAULT '0',
  `speech_details` varchar(255) DEFAULT NULL,
  `general_health_has_problem` tinyint(1) DEFAULT '0',
  `general_health_details` varchar(255) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `mental_emotional_has_problem` tinyint(1) DEFAULT '0',
  `mental_emotional_details` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `unique_idx_student_health_records_iir_id` (`iir_id`),
  KEY `idx_student_health_records_iir_id` (`iir_id`),
  CONSTRAINT `student_health_records_ibfk_1` FOREIGN KEY (`iir_id`) REFERENCES `iir_records` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=51 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `student_hobbies`
--

DROP TABLE IF EXISTS `student_hobbies`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `student_hobbies` (
  `id` int NOT NULL AUTO_INCREMENT,
  `iir_id` char(36) NOT NULL,
  `hobby_name` varchar(255) NOT NULL,
  `priority_rank` int DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_student_hobbies_iir_id` (`iir_id`),
  CONSTRAINT `student_hobbies_ibfk_1` FOREIGN KEY (`iir_id`) REFERENCES `iir_records` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=124 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `student_personal_info`
--

DROP TABLE IF EXISTS `student_personal_info`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `student_personal_info` (
  `id` int NOT NULL AUTO_INCREMENT,
  `iir_id` char(36) NOT NULL,
  `student_number` varchar(20) NOT NULL,
  `gender` enum('Male','Female') NOT NULL DEFAULT 'Male',
  `civil_status_id` int NOT NULL,
  `religion_id` int NOT NULL,
  `height_m` decimal(5,2) NOT NULL,
  `weight_kg` decimal(5,2) NOT NULL,
  `complexion` varchar(50) NOT NULL,
  `high_school_gwa` decimal(4,2) NOT NULL,
  `course_id` int NOT NULL,
  `year_level` int NOT NULL,
  `section` int NOT NULL,
  `place_of_birth` varchar(255) NOT NULL,
  `date_of_birth` date NOT NULL,
  `is_employed` tinyint(1) DEFAULT '0',
  `employer_name` varchar(255) DEFAULT NULL,
  `employer_address` varchar(255) DEFAULT NULL,
  `mobile_number` varchar(20) NOT NULL,
  `telephone_number` varchar(20) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  `status_id` int NOT NULL DEFAULT '1',
  `graduation_year` int DEFAULT NULL,
  `employer_contact_number` varchar(20) DEFAULT NULL,
  `other_religion_text` varchar(255) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `student_number` (`student_number`),
  UNIQUE KEY `unique_idx_student_personal_info_iir_id` (`iir_id`),
  UNIQUE KEY `unique_idx_student_personal_info_student_number`
    (`student_number`),
  KEY `idx_student_personal_info_iir_id` (`iir_id`),
  KEY `idx_student_personal_info_civil_status_id` (`civil_status_id`),
  KEY `idx_student_personal_info_religion_id` (`religion_id`),
  KEY `idx_student_personal_info_course_id` (`course_id`),
  KEY `idx_student_personal_info_status_id` (`status_id`),
  CONSTRAINT `fk_student_status`
    FOREIGN KEY (`status_id`) REFERENCES `student_statuses` (`id`),
  CONSTRAINT `student_personal_info_ibfk_1`
    FOREIGN KEY (`iir_id`) REFERENCES `iir_records` (`id`)
    ON DELETE CASCADE,
  CONSTRAINT `student_personal_info_ibfk_2`
    FOREIGN KEY (`religion_id`) REFERENCES `religions` (`id`),
  CONSTRAINT `student_personal_info_ibfk_3`
    FOREIGN KEY (`civil_status_id`) REFERENCES `civil_status_types` (`id`),
  CONSTRAINT `student_personal_info_ibfk_4`
    FOREIGN KEY (`course_id`) REFERENCES `courses` (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=51 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `student_related_persons`
--

DROP TABLE IF EXISTS `student_related_persons`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `student_related_persons` (
  `iir_id` char(36) NOT NULL,
  `related_person_id` int NOT NULL,
  `relationship_id` int DEFAULT NULL,
  `is_parent` tinyint(1) DEFAULT '0',
  `is_guardian` tinyint(1) DEFAULT '0',
  `is_living` tinyint(1) DEFAULT '1',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`iir_id`,`related_person_id`),
  KEY `idx_student_related_persons_iir_id` (`iir_id`),
  KEY `idx_student_related_persons_related_person_id` (`related_person_id`),
  KEY `idx_student_related_persons_relationship_id` (`relationship_id`),
  CONSTRAINT `student_related_persons_ibfk_1` FOREIGN KEY (`iir_id`) REFERENCES `iir_records` (`id`) ON DELETE CASCADE,
  CONSTRAINT `student_related_persons_ibfk_2` FOREIGN KEY (`relationship_id`) REFERENCES `student_relationship_types` (`id`),
  CONSTRAINT `student_related_persons_ibfk_3` FOREIGN KEY (`related_person_id`) REFERENCES `related_persons` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `student_relationship_types`
--

DROP TABLE IF EXISTS `student_relationship_types`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `student_relationship_types` (
  `id` int NOT NULL AUTO_INCREMENT,
  `relationship_name` varchar(100) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `relationship_name` (`relationship_name`),
  UNIQUE KEY `unique_idx_relationship_name` (`relationship_name`)
) ENGINE=InnoDB AUTO_INCREMENT=13 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `student_sibling_supports`
--

DROP TABLE IF EXISTS `student_sibling_supports`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `student_sibling_supports` (
  `family_background_id` int NOT NULL,
  `support_type_id` int NOT NULL,
  PRIMARY KEY (`family_background_id`,`support_type_id`),
  KEY `idx_student_sibling_supports_family_background_id` (`family_background_id`),
  KEY `idx_student_sibling_supports_support_type_id` (`support_type_id`),
  CONSTRAINT `student_sibling_supports_ibfk_1` FOREIGN KEY (`family_background_id`) REFERENCES `family_backgrounds` (`id`) ON DELETE CASCADE,
  CONSTRAINT `student_sibling_supports_ibfk_2` FOREIGN KEY (`support_type_id`) REFERENCES `sibling_support_types` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `student_statuses`
--

DROP TABLE IF EXISTS `student_statuses`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `student_statuses` (
  `id` int NOT NULL AUTO_INCREMENT,
  `status_name` varchar(50) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `status_name` (`status_name`)
) ENGINE=InnoDB AUTO_INCREMENT=6 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `student_subject_preferences`
--

DROP TABLE IF EXISTS `student_subject_preferences`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `student_subject_preferences` (
  `id` int NOT NULL AUTO_INCREMENT,
  `iir_id` char(36) NOT NULL,
  `subject_name` varchar(100) NOT NULL,
  `is_favorite` tinyint(1) DEFAULT '0',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `iir_id` (`iir_id`,`subject_name`),
  UNIQUE KEY `unique_idx_student_subject_preferences_iir_id_subject_name` (`iir_id`,`subject_name`),
  KEY `idx_student_subject_preferences_iir_id` (`iir_id`),
  CONSTRAINT `student_subject_preferences_ibfk_1` FOREIGN KEY (`iir_id`) REFERENCES `iir_records` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=165 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `student_support_types`
--

DROP TABLE IF EXISTS `student_support_types`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `student_support_types` (
  `id` int NOT NULL AUTO_INCREMENT,
  `support_type_name` varchar(100) NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `support_type_name` (`support_type_name`),
  UNIQUE KEY `unique_idx_student_support_type_name` (`support_type_name`)
) ENGINE=InnoDB AUTO_INCREMENT=7 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `system_logs`
--

DROP TABLE IF EXISTS `system_logs`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `system_logs` (
  `id` int NOT NULL AUTO_INCREMENT,
  `level` enum('INFO','WARNING','ERROR','CRITICAL') NOT NULL,
  `category` enum('SECURITY','SYSTEM','AUDIT','CONSENT') NOT NULL,
  `action` varchar(100) NOT NULL,
  `message` text NOT NULL,
  `user_id` char(36) DEFAULT NULL,
  `target_id` char(36) DEFAULT NULL,
  `target_type` varchar(50) DEFAULT NULL,
  `user_email` varchar(100) DEFAULT NULL,
  `target_email` varchar(100) DEFAULT NULL,
  `ip_address` varchar(45) DEFAULT NULL,
  `user_agent` varchar(255) DEFAULT NULL,
  `metadata` json DEFAULT NULL,
  `trace_id` char(36) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_system_logs_action` (`action`),
  KEY `idx_system_logs_user_email` (`user_email`),
  KEY `idx_system_logs_target_email` (`target_email`),
  KEY `idx_system_logs_created_at` (`created_at`),
  KEY `idx_system_logs_category_created` (`created_at`),
  KEY `idx_system_logs_user_id` (`user_id`)
) ENGINE=InnoDB AUTO_INCREMENT=93 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `test_results`
--

DROP TABLE IF EXISTS `test_results`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `test_results` (
  `id` int NOT NULL AUTO_INCREMENT,
  `iir_id` char(36) NOT NULL,
  `test_date` date DEFAULT NULL,
  `test_name` varchar(255) DEFAULT NULL,
  `raw_score` varchar(50) DEFAULT NULL,
  `percentile` varchar(50) DEFAULT NULL,
  `description` varchar(255) DEFAULT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  KEY `idx_test_results_iir_id` (`iir_id`),
  CONSTRAINT `test_results_ibfk_1` FOREIGN KEY (`iir_id`) REFERENCES `iir_records` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB AUTO_INCREMENT=101 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `time_slots`
--

DROP TABLE IF EXISTS `time_slots`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `time_slots` (
  `id` int NOT NULL AUTO_INCREMENT,
  `time` time NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `time` (`time`),
  UNIQUE KEY `unique_idx_time_slot_time` (`time`)
) ENGINE=InnoDB AUTO_INCREMENT=9 DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `user_roles`
--

DROP TABLE IF EXISTS `user_roles`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `user_roles` (
  `user_id` char(36) NOT NULL,
  `role_id` int NOT NULL,
  `assigned_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `assigned_by` char(36) DEFAULT NULL,
  `reason` text,
  `reference_id` varchar(100) DEFAULT NULL,
  PRIMARY KEY (`user_id`,`role_id`),
  KEY `fk_user_roles_role` (`role_id`),
  KEY `fk_user_roles_admin` (`assigned_by`),
  CONSTRAINT `fk_user_roles_admin` FOREIGN KEY (`assigned_by`) REFERENCES `users` (`id`) ON DELETE SET NULL,
  CONSTRAINT `fk_user_roles_role` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`) ON DELETE CASCADE,
  CONSTRAINT `fk_user_roles_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Table structure for table `users`
--

DROP TABLE IF EXISTS `users`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `users` (
  `id` char(36) NOT NULL,
  `email` varchar(100) NOT NULL,
  `first_name` varchar(100) NOT NULL,
  `middle_name` varchar(100) DEFAULT NULL,
  `last_name` varchar(100) NOT NULL,
  `suffix_name` varchar(50) DEFAULT NULL,
  `password_hash` varchar(255) DEFAULT NULL,
  `auth_type` varchar(20) NOT NULL DEFAULT 'native',
  `is_active` tinyint(1) DEFAULT '1',
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_users_email_auth_type` (`email`,`auth_type`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Temporary view structure for view `v_related_persons`
--

DROP TABLE IF EXISTS `v_related_persons`;
/*!50001 DROP VIEW IF EXISTS `v_related_persons`*/;
SET @saved_cs_client     = @@character_set_client;
/*!50503 SET character_set_client = utf8mb4 */;
/*!50001 CREATE VIEW `v_related_persons` AS SELECT 
 1 AS `id`,
 1 AS `last_name`,
 1 AS `first_name`,
 1 AS `middle_name`,
 1 AS `suffix_name`,
 1 AS `date_of_birth`,
 1 AS `educational_attainment_id`,
 1 AS `educational_attainment_name`,
 1 AS `occupation`,
 1 AS `employer_name`,
 1 AS `employer_address`,
 1 AS `iir_id`,
 1 AS `relationship_id`,
 1 AS `relationship_name`,
 1 AS `is_parent`,
 1 AS `is_guardian`,
 1 AS `is_living`*/;
SET character_set_client = @saved_cs_client;

--
-- Temporary view structure for view `v_student_basic_info`
--

DROP TABLE IF EXISTS `v_student_basic_info`;
/*!50001 DROP VIEW IF EXISTS `v_student_basic_info`*/;
SET @saved_cs_client     = @@character_set_client;
/*!50503 SET character_set_client = utf8mb4 */;
/*!50001 CREATE VIEW `v_student_basic_info` AS SELECT 
 1 AS `iir_id`,
 1 AS `user_id`,
 1 AS `email`,
 1 AS `first_name`,
 1 AS `middle_name`,
 1 AS `last_name`,
 1 AS `suffix_name`*/;
SET character_set_client = @saved_cs_client;

--
-- Temporary view structure for view `v_student_cors_files`
--

DROP TABLE IF EXISTS `v_student_cors_files`;
/*!50001 DROP VIEW IF EXISTS `v_student_cors_files`*/;
SET @saved_cs_client     = @@character_set_client;
/*!50503 SET character_set_client = utf8mb4 */;
/*!50001 CREATE VIEW `v_student_cors_files` AS SELECT 
 1 AS `student_id`,
 1 AS `file_url`*/;
SET character_set_client = @saved_cs_client;

--
-- Temporary view structure for view `v_student_current_cors`
--

DROP TABLE IF EXISTS `v_student_current_cors`;
/*!50001 DROP VIEW IF EXISTS `v_student_current_cors`*/;
SET @saved_cs_client     = @@character_set_client;
/*!50503 SET character_set_client = utf8mb4 */;
/*!50001 CREATE VIEW `v_student_current_cors` AS SELECT 
 1 AS `file_id`,
 1 AS `student_id`,
 1 AS `student_number`,
 1 AS `course_desc`,
 1 AS `course_code`,
 1 AS `year_level`,
 1 AS `section`,
 1 AS `campus`,
 1 AS `year_start`,
 1 AS `year_end`,
 1 AS `term`,
 1 AS `valid_from`,
 1 AS `valid_until`*/;
SET character_set_client = @saved_cs_client;

--
-- Temporary view structure for view `v_student_finances`
--

DROP TABLE IF EXISTS `v_student_finances`;
/*!50001 DROP VIEW IF EXISTS `v_student_finances`*/;
SET @saved_cs_client     = @@character_set_client;
/*!50503 SET character_set_client = utf8mb4 */;
/*!50001 CREATE VIEW `v_student_finances` AS SELECT 
 1 AS `id`,
 1 AS `iir_id`,
 1 AS `income_range_id`,
 1 AS `income_range_text`,
 1 AS `other_income`,
 1 AS `weekly_allowance`*/;
SET character_set_client = @saved_cs_client;

--
-- Temporary view structure for view `v_student_financial_supports`
--

DROP TABLE IF EXISTS `v_student_financial_supports`;
/*!50001 DROP VIEW IF EXISTS `v_student_financial_supports`*/;
SET @saved_cs_client     = @@character_set_client;
/*!50503 SET character_set_client = utf8mb4 */;
/*!50001 CREATE VIEW `v_student_financial_supports` AS SELECT 
 1 AS `sf_id`,
 1 AS `id`,
 1 AS `support_type_name`*/;
SET character_set_client = @saved_cs_client;

--
-- Temporary view structure for view `v_student_personal_info`
--

DROP TABLE IF EXISTS `v_student_personal_info`;
/*!50001 DROP VIEW IF EXISTS `v_student_personal_info`*/;
SET @saved_cs_client     = @@character_set_client;
/*!50503 SET character_set_client = utf8mb4 */;
/*!50001 CREATE VIEW `v_student_personal_info` AS SELECT 
 1 AS `id`,
 1 AS `iir_id`,
 1 AS `student_number`,
 1 AS `gender_id`,
 1 AS `gender_name`,
 1 AS `civil_status_id`,
 1 AS `civil_status_name`,
 1 AS `religion_id`,
 1 AS `religion_name`,
 1 AS `other_religion_text`,
 1 AS `height_m`,
 1 AS `weight_kg`,
 1 AS `complexion`,
 1 AS `high_school_gwa`,
 1 AS `course_id`,
 1 AS `course_code`,
 1 AS `course_name`,
 1 AS `year_level`,
 1 AS `section`,
 1 AS `place_of_birth`,
 1 AS `date_of_birth`,
 1 AS `is_employed`,
 1 AS `employer_name`,
 1 AS `employer_address`,
 1 AS `mobile_number`,
 1 AS `telephone_number`,
 1 AS `employer_contact_number`,
 1 AS `two_by_two_photo_data_url`,
 1 AS `status_id`,
 1 AS `status_name`,
 1 AS `graduation_year`,
 1 AS `emergency_id`,
 1 AS `emergency_first_name`,
 1 AS `emergency_middle_name`,
 1 AS `emergency_last_name`,
 1 AS `emergency_contact_number`,
 1 AS `emergency_relationship_id`,
 1 AS `emergency_relationship_name`,
 1 AS `emergency_address_id`*/;
SET character_set_client = @saved_cs_client;

--
-- Temporary view structure for view `v_student_profiles`
--

DROP TABLE IF EXISTS `v_student_profiles`;
/*!50001 DROP VIEW IF EXISTS `v_student_profiles`*/;
SET @saved_cs_client     = @@character_set_client;
/*!50503 SET character_set_client = utf8mb4 */;
/*!50001 CREATE VIEW `v_student_profiles` AS SELECT 
 1 AS `iir_id`,
 1 AS `user_id`,
 1 AS `first_name`,
 1 AS `middle_name`,
 1 AS `last_name`,
 1 AS `suffix_name`,
 1 AS `email`,
 1 AS `student_number`,
 1 AS `gender_id`,
 1 AS `course_id`,
 1 AS `section`,
 1 AS `year_level`,
 1 AS `status_id`,
 1 AS `status_name`,
 1 AS `gender_name`,
 1 AS `profile_picture`,
 1 AS `course_code`,
 1 AS `course_name`,
 1 AS `created_at`,
 1 AS `updated_at`*/;
SET character_set_client = @saved_cs_client;

--
-- Table structure for table `whitelists`
--

DROP TABLE IF EXISTS `whitelists`;
/*!40101 SET @saved_cs_client     = @@character_set_client */;
/*!50503 SET character_set_client = utf8mb4 */;
CREATE TABLE `whitelists` (
  `email` varchar(100) NOT NULL,
  `role_id` int NOT NULL,
  `created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`email`,`role_id`),
  KEY `whitelists_ibfk_1` (`role_id`),
  CONSTRAINT `whitelists_ibfk_1` FOREIGN KEY (`role_id`) REFERENCES `roles` (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;
/*!40101 SET character_set_client = @saved_cs_client */;

--
-- Final view structure for view `v_related_persons`
--

/*!50001 DROP VIEW IF EXISTS `v_related_persons`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_general_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`mysqladmin`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `v_related_persons` AS select `rp`.`id` AS `id`,`rp`.`last_name` AS `last_name`,`rp`.`first_name` AS `first_name`,`rp`.`middle_name` AS `middle_name`,`rp`.`suffix_name` AS `suffix_name`,`rp`.`date_of_birth` AS `date_of_birth`,coalesce(`rp`.`educational_attainment_id`,0) AS `educational_attainment_id`,coalesce(`ea`.`name`,'') AS `educational_attainment_name`,coalesce(`rp`.`occupation`,'') AS `occupation`,coalesce(`rp`.`employer_name`,'') AS `employer_name`,coalesce(`rp`.`employer_address`,'') AS `employer_address`,`srp`.`iir_id` AS `iir_id`,`srp`.`relationship_id` AS `relationship_id`,coalesce(`ert`.`relationship_name`,'') AS `relationship_name`,`srp`.`is_parent` AS `is_parent`,`srp`.`is_guardian` AS `is_guardian`,`srp`.`is_living` AS `is_living` from (((`student_related_persons` `srp` join `related_persons` `rp` on((`srp`.`related_person_id` = `rp`.`id`))) left join `educational_attainments` `ea` on((`rp`.`educational_attainment_id` = `ea`.`id`))) left join `student_relationship_types` `ert` on((`srp`.`relationship_id` = `ert`.`id`))) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `v_student_basic_info`
--

/*!50001 DROP VIEW IF EXISTS `v_student_basic_info`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_general_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`mysqladmin`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `v_student_basic_info` AS select `iir`.`id` AS `iir_id`,`u`.`id` AS `user_id`,`u`.`email` AS `email`,`u`.`first_name` AS `first_name`,`u`.`middle_name` AS `middle_name`,`u`.`last_name` AS `last_name`,`u`.`suffix_name` AS `suffix_name` from (`users` `u` join `iir_records` `iir` on((`u`.`id` = `iir`.`user_id`))) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `v_student_cors_files`
--

/*!50001 DROP VIEW IF EXISTS `v_student_cors_files`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_general_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`mysqladmin`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `v_student_cors_files` AS select `sc`.`student_id` AS `student_id`,`f`.`file_url` AS `file_url` from (`student_cors` `sc` join `files` `f` on((`f`.`id` = `sc`.`file_id`))) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `v_student_current_cors`
--

/*!50001 DROP VIEW IF EXISTS `v_student_current_cors`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_general_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`mysqladmin`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `v_student_current_cors` AS select `sc`.`file_id` AS `file_id`,`sc`.`student_id` AS `student_id`,`sc`.`student_number` AS `student_number`,`sc`.`course_desc` AS `course_desc`,`sc`.`course_code` AS `course_code`,`sc`.`year_level` AS `year_level`,`sc`.`section` AS `section`,`sc`.`campus` AS `campus`,`sc`.`year_start` AS `year_start`,`sc`.`year_end` AS `year_end`,`sc`.`term` AS `term`,`sc`.`valid_from` AS `valid_from`,`sc`.`valid_until` AS `valid_until` from (`student_cors` `sc` join `academic_settings` `ac` on((`ac`.`id` = 1))) where ((`sc`.`year_start` = `ac`.`current_year_start`) and (`sc`.`term` = `ac`.`current_term`) and (`sc`.`valid_from` is not null) and (`sc`.`valid_until` is not null)) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `v_student_finances`
--

/*!50001 DROP VIEW IF EXISTS `v_student_finances`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_general_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`mysqladmin`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `v_student_finances` AS select `sf`.`id` AS `id`,`sf`.`iir_id` AS `iir_id`,`sf`.`monthly_family_income_range_id` AS `income_range_id`,coalesce(`ir`.`range_text`,'') AS `income_range_text`,`sf`.`other_income_details` AS `other_income`,`sf`.`weekly_allowance` AS `weekly_allowance` from (`student_finances` `sf` left join `income_ranges` `ir` on((`sf`.`monthly_family_income_range_id` = `ir`.`id`))) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `v_student_financial_supports`
--

/*!50001 DROP VIEW IF EXISTS `v_student_financial_supports`*/;
/*!50001 SET @saved_cs_client          = @@character_set_client */;
/*!50001 SET @saved_cs_results         = @@character_set_results */;
/*!50001 SET @saved_col_connection     = @@collation_connection */;
/*!50001 SET character_set_client      = utf8mb4 */;
/*!50001 SET character_set_results     = utf8mb4 */;
/*!50001 SET collation_connection      = utf8mb4_general_ci */;
/*!50001 CREATE ALGORITHM=UNDEFINED */
/*!50013 DEFINER=`mysqladmin`@`%` SQL SECURITY DEFINER */
/*!50001 VIEW `v_student_financial_supports` AS select `sfs`.`sf_id` AS `sf_id`,`sst`.`id` AS `id`,`sst`.`support_type_name` AS `support_type_name` from (`student_financial_supports` `sfs` join `student_support_types` `sst` on((`sfs`.`support_type_id` = `sst`.`id`))) */;
/*!50001 SET character_set_client      = @saved_cs_client */;
/*!50001 SET character_set_results     = @saved_cs_results */;
/*!50001 SET collation_connection      = @saved_col_connection */;

--
-- Final view structure for view `v_student_personal_info`
--

DROP VIEW IF EXISTS `v_student_personal_info`;
CREATE OR REPLACE VIEW `v_student_personal_info` AS
SELECT
  spi.id,
  spi.iir_id,
  spi.student_number,
  CASE WHEN spi.gender = 'Male' THEN 1 ELSE 2 END AS gender_id,
  spi.gender AS gender_name,
  spi.civil_status_id,
  COALESCE(cst.status_name, '') AS civil_status_name,
  spi.religion_id,
  COALESCE(rel.religion_name, '') AS religion_name,
  spi.other_religion_text,
  spi.height_m,
  spi.weight_kg,
  spi.complexion,
  spi.high_school_gwa,
  spi.course_id,
  COALESCE(c.code, '') AS course_code,
  COALESCE(c.course_name, '') AS course_name,
  spi.year_level,
  spi.section,
  spi.place_of_birth,
  spi.date_of_birth,
  spi.is_employed,
  spi.employer_name,
  spi.employer_address,
  spi.mobile_number,
  spi.telephone_number,
  spi.employer_contact_number,
  COALESCE(pf.file_url, '') AS two_by_two_photo_data_url,
  spi.status_id,
  COALESCE(ss.status_name, '') AS status_name,
  spi.graduation_year,
  COALESCE(ec.id, 0) AS emergency_id,
  COALESCE(ec.first_name, '') AS emergency_first_name,
  COALESCE(ec.middle_name, '') AS emergency_middle_name,
  COALESCE(ec.last_name, '') AS emergency_last_name,
  COALESCE(ec.contact_number, '') AS emergency_contact_number,
  COALESCE(ec.relationship_id, 0) AS emergency_relationship_id,
  COALESCE(ert.relationship_name, '') AS emergency_relationship_name,
  COALESCE(ec.address_id, 0) AS emergency_address_id
FROM student_personal_info spi
JOIN iir_records iir ON iir.id = spi.iir_id
LEFT JOIN profile_pictures pp ON pp.user_id = iir.user_id
LEFT JOIN files pf ON pf.id = pp.file_id
LEFT JOIN civil_status_types cst ON spi.civil_status_id = cst.id
LEFT JOIN religions rel ON spi.religion_id = rel.id
LEFT JOIN courses c ON spi.course_id = c.id
LEFT JOIN student_statuses ss ON spi.status_id = ss.id
LEFT JOIN emergency_contacts ec ON spi.iir_id = ec.iir_id
LEFT JOIN student_relationship_types ert ON ec.relationship_id = ert.id;

--
-- Final view structure for view `v_student_profiles`
--

DROP VIEW IF EXISTS `v_student_profiles`;
CREATE OR REPLACE VIEW `v_student_profiles` AS
SELECT
  iir.id AS iir_id,
  iir.user_id,
  u.first_name,
  u.middle_name,
  u.last_name,
  u.suffix_name,
  u.email,
  spi.student_number,
  CASE WHEN spi.gender = 'Male' THEN 1 ELSE 2 END AS gender_id,
  spi.course_id,
  spi.section,
  spi.year_level,
  spi.status_id,
  COALESCE(ss.status_name, '') AS status_name,
  spi.gender AS gender_name,
  COALESCE(pf.file_url, '') AS profile_picture,
  COALESCE(c.code, '') AS course_code,
  COALESCE(c.course_name, '') AS course_name,
  iir.created_at,
  iir.updated_at
FROM iir_records iir
JOIN users u ON iir.user_id = u.id
JOIN student_personal_info spi ON iir.id = spi.iir_id
LEFT JOIN profile_pictures pp ON pp.user_id = iir.user_id
LEFT JOIN files pf ON pf.id = pp.file_id
LEFT JOIN student_statuses ss ON spi.status_id = ss.id
LEFT JOIN courses c ON spi.course_id = c.id
WHERE iir.is_submitted = TRUE;
/*!40103 SET TIME_ZONE=@OLD_TIME_ZONE */;

/*!40101 SET SQL_MODE=@OLD_SQL_MODE */;
/*!40014 SET FOREIGN_KEY_CHECKS=@OLD_FOREIGN_KEY_CHECKS */;
/*!40014 SET UNIQUE_CHECKS=@OLD_UNIQUE_CHECKS */;
/*!40101 SET CHARACTER_SET_CLIENT=@OLD_CHARACTER_SET_CLIENT */;
/*!40101 SET CHARACTER_SET_RESULTS=@OLD_CHARACTER_SET_RESULTS */;
/*!40101 SET COLLATION_CONNECTION=@OLD_COLLATION_CONNECTION */;
/*!40111 SET SQL_NOTES=@OLD_SQL_NOTES */;

-- Dump completed on 2026-06-15 15:44:27
