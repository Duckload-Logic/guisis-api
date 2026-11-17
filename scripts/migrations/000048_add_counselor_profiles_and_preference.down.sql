ALTER TABLE `appointments`
DROP FOREIGN KEY `fk_appointments_preferred_counselor`,
DROP COLUMN `counselor_preference`,
DROP COLUMN `preferred_counselor_id`;

DROP TABLE IF EXISTS `counselor_profiles`;
