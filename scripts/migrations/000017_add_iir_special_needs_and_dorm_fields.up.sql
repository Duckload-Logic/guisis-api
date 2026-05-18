ALTER TABLE student_personal_info
ADD COLUMN employer_contact_number VARCHAR(20) DEFAULT NULL,
ADD COLUMN living_in_dorm TINYINT(1) DEFAULT 0,
ADD COLUMN dorm_address VARCHAR(255) DEFAULT NULL,
ADD COLUMN landlord_name VARCHAR(255) DEFAULT NULL,
ADD COLUMN landlord_contact_number VARCHAR(20) DEFAULT NULL;

ALTER TABLE student_health_records
ADD COLUMN mental_emotional_has_problem TINYINT(1) DEFAULT 0,
ADD COLUMN mental_emotional_details VARCHAR(255) DEFAULT NULL;
