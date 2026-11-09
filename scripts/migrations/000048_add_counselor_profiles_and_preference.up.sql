CREATE TABLE IF NOT EXISTS `counselor_profiles` (
  `id` INT AUTO_INCREMENT PRIMARY KEY,
  `user_id` CHAR(36) NOT NULL UNIQUE,
  `title` VARCHAR(20) NOT NULL DEFAULT '',
  `gender` ENUM('FEMALE', 'MALE', 'OTHER') NOT NULL DEFAULT 'FEMALE',
  `position` VARCHAR(100) NOT NULL DEFAULT 'Guidance Counselor',
  `is_active` TINYINT(1) NOT NULL DEFAULT 1,
  `created_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  `updated_at` TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  CONSTRAINT `fk_counselor_profiles_user` 
    FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;

ALTER TABLE `appointments`
ADD COLUMN `preferred_counselor_id` CHAR(36) NULL AFTER `urgency_rationale`,
ADD COLUMN `counselor_preference` VARCHAR(100) NOT NULL DEFAULT 'NO PREFERENCE' AFTER `preferred_counselor_id`,
ADD CONSTRAINT `fk_appointments_preferred_counselor` 
  FOREIGN KEY (`preferred_counselor_id`) REFERENCES `users` (`id`) ON DELETE SET NULL;
