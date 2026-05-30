-- ============================================================================
-- CORE IIR DOMAIN
-- ============================================================================

CREATE TABLE iir_records (
    id CHAR(36) NOT NULL PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    is_submitted TINYINT(1) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_iir_records_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE UNIQUE INDEX idx_iir_records_user_unique ON iir_records(user_id ASC);
CREATE INDEX idx_iir_records_user_id ON iir_records(user_id ASC);

CREATE TABLE educational_backgrounds (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    iir_id CHAR(36) NOT NULL,
    nature_of_schooling ENUM('Continuous', 'Interrupted') NOT NULL,
    interrupted_details VARCHAR(255) DEFAULT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT educational_backgrounds_ibfk_1 FOREIGN KEY (iir_id) REFERENCES iir_records(id) ON DELETE CASCADE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE INDEX idx_educational_backgrounds_iir_id ON educational_backgrounds(iir_id ASC);

CREATE TABLE school_details (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    eb_id INT NOT NULL,
    educational_level_id INT NOT NULL,
    school_name VARCHAR(255) NOT NULL,
    school_address VARCHAR(255) NOT NULL,
    school_type ENUM('Private', 'Public') NOT NULL,
    year_started SMALLINT NOT NULL,
    year_completed SMALLINT NOT NULL,
    awards TEXT DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT school_details_ibfk_1 FOREIGN KEY (eb_id) REFERENCES educational_backgrounds(id) ON DELETE CASCADE,
    CONSTRAINT school_details_ibfk_2 FOREIGN KEY (educational_level_id) REFERENCES educational_levels(id)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE INDEX idx_school_details_eb_id ON school_details(eb_id ASC);
CREATE INDEX idx_school_details_educational_level_id ON school_details(educational_level_id ASC);

CREATE TABLE family_backgrounds (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    iir_id CHAR(36) NOT NULL,
    parental_status_id INT NOT NULL,
    parental_status_details VARCHAR(255) DEFAULT NULL,
    brothers INT NOT NULL,
    sisters INT NOT NULL,
    employed_siblings INT NOT NULL,
    ordinal_position INT NOT NULL,
    have_quiet_place_to_study TINYINT(1) NOT NULL,
    is_sharing_room TINYINT(1) NOT NULL,
    room_sharing_details VARCHAR(255) DEFAULT NULL,
    nature_of_residence_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT family_backgrounds_ibfk_1 FOREIGN KEY (iir_id) REFERENCES iir_records(id) ON DELETE CASCADE,
    CONSTRAINT family_backgrounds_ibfk_2 FOREIGN KEY (nature_of_residence_id) REFERENCES nature_of_residence_types(id),
    CONSTRAINT family_backgrounds_ibfk_3 FOREIGN KEY (parental_status_id) REFERENCES parental_status_types(id)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE UNIQUE INDEX iir_id ON family_backgrounds(iir_id ASC);
CREATE INDEX idx_family_backgrounds_iir_id ON family_backgrounds(iir_id ASC);
CREATE INDEX idx_family_backgrounds_parental_status_id ON family_backgrounds(parental_status_id ASC);
CREATE INDEX idx_family_backgrounds_nature_of_residence_id ON family_backgrounds(nature_of_residence_id ASC);

CREATE TABLE student_sibling_supports (
    family_background_id INT NOT NULL,
    support_type_id INT NOT NULL,
    PRIMARY KEY (family_background_id, support_type_id),
    CONSTRAINT student_sibling_supports_ibfk_1 FOREIGN KEY (family_background_id) REFERENCES family_backgrounds(id) ON DELETE CASCADE,
    CONSTRAINT student_sibling_supports_ibfk_2 FOREIGN KEY (support_type_id) REFERENCES sibling_support_types(id) ON DELETE CASCADE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE INDEX idx_student_sibling_supports_family_background_id ON student_sibling_supports(family_background_id ASC);
CREATE INDEX idx_student_sibling_supports_support_type_id ON student_sibling_supports(support_type_id ASC);

CREATE TABLE emergency_contacts (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    iir_id CHAR(36) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(100) DEFAULT NULL,
    last_name VARCHAR(100) NOT NULL,
    suffix_name VARCHAR(50) DEFAULT NULL,
    contact_number VARCHAR(20) NOT NULL,
    relationship_id INT NOT NULL,
    address_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT emergency_contacts_ibfk_1 FOREIGN KEY (iir_id) REFERENCES iir_records(id) ON DELETE CASCADE,
    CONSTRAINT emergency_contacts_ibfk_2 FOREIGN KEY (relationship_id) REFERENCES student_relationship_types(id),
    CONSTRAINT emergency_contacts_ibfk_3 FOREIGN KEY (address_id) REFERENCES addresses(id)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE INDEX idx_emergency_contacts_iir_id ON emergency_contacts(iir_id ASC);
CREATE INDEX idx_emergency_contacts_relationship_id ON emergency_contacts(relationship_id ASC);
CREATE INDEX idx_emergency_contacts_address_id ON emergency_contacts(address_id ASC);

CREATE TABLE student_personal_info (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    iir_id CHAR(36) NOT NULL,
    student_number VARCHAR(20) NOT NULL UNIQUE,
    gender_id INT NOT NULL,
    civil_status_id INT NOT NULL,
    religion_id INT NOT NULL,
    height_m DECIMAL(5,2) NOT NULL,
    weight_kg DECIMAL(5,2) NOT NULL,
    complexion VARCHAR(50) NOT NULL,
    high_school_gwa DECIMAL(4,2) NOT NULL,
    course_id INT NOT NULL,
    year_level INT NOT NULL,
    section INT NOT NULL,
    place_of_birth VARCHAR(255) NOT NULL,
    date_of_birth DATE NOT NULL,
    is_employed TINYINT(1) DEFAULT 0,
    employer_name VARCHAR(255) DEFAULT NULL,
    employer_address VARCHAR(255) DEFAULT NULL,
    mobile_number VARCHAR(20) NOT NULL,
    telephone_number VARCHAR(20) DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT student_personal_info_ibfk_1 FOREIGN KEY (iir_id) REFERENCES iir_records(id) ON DELETE CASCADE,
    CONSTRAINT student_personal_info_ibfk_2 FOREIGN KEY (religion_id) REFERENCES religions(id),
    CONSTRAINT student_personal_info_ibfk_3 FOREIGN KEY (civil_status_id) REFERENCES civil_status_types(id),
    CONSTRAINT student_personal_info_ibfk_4 FOREIGN KEY (course_id) REFERENCES courses(id),
    CONSTRAINT student_personal_info_ibfk_5 FOREIGN KEY (gender_id) REFERENCES genders(id)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE UNIQUE INDEX unique_idx_student_personal_info_iir_id ON student_personal_info(iir_id ASC);
CREATE UNIQUE INDEX unique_idx_student_personal_info_student_number ON student_personal_info(student_number ASC);
CREATE INDEX idx_student_personal_info_iir_id ON student_personal_info(iir_id ASC);
CREATE INDEX idx_student_personal_info_gender_id ON student_personal_info(gender_id ASC);
CREATE INDEX idx_student_personal_info_civil_status_id ON student_personal_info(civil_status_id ASC);
CREATE INDEX idx_student_personal_info_religion_id ON student_personal_info(religion_id ASC);
CREATE INDEX idx_student_personal_info_course_id ON student_personal_info(course_id ASC);

CREATE TABLE student_addresses (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    iir_id CHAR(36) NOT NULL,
    address_id INT NOT NULL,
    address_type ENUM('Residential', 'Provincial') NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT student_addresses_ibfk_1 FOREIGN KEY (iir_id) REFERENCES iir_records(id) ON DELETE CASCADE,
    CONSTRAINT student_addresses_ibfk_2 FOREIGN KEY (address_id) REFERENCES addresses(id) ON DELETE CASCADE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE INDEX idx_student_addresses_iir_id ON student_addresses(iir_id ASC);
CREATE INDEX idx_student_addresses_address_id ON student_addresses(address_id ASC);

CREATE TABLE related_persons (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(100) DEFAULT NULL,
    last_name VARCHAR(100) NOT NULL,
    suffix_name VARCHAR(50) DEFAULT NULL,
    educational_level VARCHAR(100) NOT NULL,
    date_of_birth DATE NOT NULL,
    occupation VARCHAR(100) DEFAULT NULL,
    employer_name VARCHAR(150) DEFAULT NULL,
    employer_address VARCHAR(255) DEFAULT NULL,
    contact_number VARCHAR(20) DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE TABLE student_related_persons (
    iir_id CHAR(36) NOT NULL,
    related_person_id INT NOT NULL,
    relationship_id INT DEFAULT NULL,
    is_parent TINYINT(1) DEFAULT 0,
    is_guardian TINYINT(1) DEFAULT 0,
    is_living TINYINT(1) DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    PRIMARY KEY (iir_id, related_person_id),
    CONSTRAINT student_related_persons_ibfk_1 FOREIGN KEY (iir_id) REFERENCES iir_records(id) ON DELETE CASCADE,
    CONSTRAINT student_related_persons_ibfk_2 FOREIGN KEY (relationship_id) REFERENCES student_relationship_types(id),
    CONSTRAINT student_related_persons_ibfk_3 FOREIGN KEY (related_person_id) REFERENCES related_persons(id) ON DELETE CASCADE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE INDEX idx_student_related_persons_iir_id ON student_related_persons(iir_id ASC);
CREATE INDEX idx_student_related_persons_related_person_id ON student_related_persons(related_person_id ASC);
CREATE INDEX idx_student_related_persons_relationship_id ON student_related_persons(relationship_id ASC);

CREATE TABLE student_activities (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    iir_id CHAR(36) NOT NULL,
    option_id INT DEFAULT NULL,
    other_specification VARCHAR(255) DEFAULT NULL,
    role ENUM('Officer', 'Member', 'Other') DEFAULT 'Member',
    role_specification VARCHAR(255) DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT student_activities_ibfk_1 FOREIGN KEY (iir_id) REFERENCES iir_records(id) ON DELETE CASCADE,
    CONSTRAINT student_activities_ibfk_2 FOREIGN KEY (option_id) REFERENCES activity_options(id)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE INDEX idx_student_activities_iir_id ON student_activities(iir_id ASC);
CREATE INDEX idx_student_activities_option_id ON student_activities(option_id ASC);

CREATE TABLE student_consultations (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    iir_id CHAR(36) NOT NULL,
    professional_type ENUM('Psychiatrist', 'Psychologist', 'Counselor') NOT NULL,
    has_consulted TINYINT(1) DEFAULT 0,
    when_date VARCHAR(100) DEFAULT NULL,
    for_what TEXT DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT student_consultations_ibfk_1 FOREIGN KEY (iir_id) REFERENCES iir_records(id) ON DELETE CASCADE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE INDEX idx_student_consultations_iir_id ON student_consultations(iir_id ASC);

CREATE TABLE student_finances (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    iir_id CHAR(36) NOT NULL,
    monthly_family_income_range_id INT DEFAULT NULL,
    other_income_details VARCHAR(50) DEFAULT NULL,
    weekly_allowance DECIMAL(10,2) DEFAULT 0.00,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT student_finances_ibfk_1 FOREIGN KEY (iir_id) REFERENCES iir_records(id) ON DELETE CASCADE,
    CONSTRAINT student_finances_ibfk_2 FOREIGN KEY (monthly_family_income_range_id) REFERENCES income_ranges(id)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE UNIQUE INDEX unique_idx_student_finances_iir_id ON student_finances(iir_id ASC);
CREATE INDEX idx_student_finances_iir_id ON student_finances(iir_id ASC);
CREATE INDEX idx_student_finances_income_range_id ON student_finances(monthly_family_income_range_id ASC);

CREATE TABLE student_financial_supports (
    sf_id INT NOT NULL,
    support_type_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT student_financial_supports_ibfk_1 FOREIGN KEY (support_type_id) REFERENCES student_support_types(id),
    CONSTRAINT student_financial_supports_ibfk_2 FOREIGN KEY (sf_id) REFERENCES student_finances(id) ON DELETE CASCADE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE INDEX idx_student_financial_supports_sf_id ON student_financial_supports(sf_id ASC);
CREATE INDEX idx_student_financial_supports_support_type_id ON student_financial_supports(support_type_id ASC);

CREATE TABLE student_health_records (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    iir_id CHAR(36) NOT NULL,
    vision_has_problem TINYINT(1) DEFAULT 0,
    vision_details VARCHAR(255) DEFAULT NULL,
    hearing_has_problem TINYINT(1) DEFAULT 0,
    hearing_details VARCHAR(255) DEFAULT NULL,
    speech_has_problem TINYINT(1) DEFAULT 0,
    speech_details VARCHAR(255) DEFAULT NULL,
    general_health_has_problem TINYINT(1) DEFAULT 0,
    general_health_details VARCHAR(255) DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT student_health_records_ibfk_1 FOREIGN KEY (iir_id) REFERENCES iir_records(id) ON DELETE CASCADE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE UNIQUE INDEX unique_idx_student_health_records_iir_id ON student_health_records(iir_id ASC);
CREATE INDEX idx_student_health_records_iir_id ON student_health_records(iir_id ASC);

CREATE TABLE student_hobbies (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    iir_id CHAR(36) NOT NULL,
    hobby_name VARCHAR(255) NOT NULL,
    priority_rank INT DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT student_hobbies_ibfk_1 FOREIGN KEY (iir_id) REFERENCES iir_records(id) ON DELETE CASCADE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE INDEX idx_student_hobbies_iir_id ON student_hobbies(iir_id ASC);

CREATE TABLE student_subject_preferences (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    iir_id CHAR(36) NOT NULL,
    subject_name VARCHAR(100) NOT NULL,
    is_favorite TINYINT(1) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT student_subject_preferences_ibfk_1 FOREIGN KEY (iir_id) REFERENCES iir_records(id) ON DELETE CASCADE,
    UNIQUE KEY (iir_id, subject_name)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE UNIQUE INDEX unique_idx_student_subject_preferences_iir_id_subject_name ON student_subject_preferences(iir_id ASC, subject_name ASC);
CREATE INDEX idx_student_subject_preferences_iir_id ON student_subject_preferences(iir_id ASC);

CREATE TABLE test_results (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    iir_id CHAR(36) NOT NULL,
    test_date DATE DEFAULT NULL,
    test_name VARCHAR(255) DEFAULT NULL,
    raw_score VARCHAR(50) DEFAULT NULL,
    percentile VARCHAR(50) DEFAULT NULL,
    description VARCHAR(255) DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT test_results_ibfk_1 FOREIGN KEY (iir_id) REFERENCES iir_records(id) ON DELETE CASCADE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE INDEX idx_test_results_iir_id ON test_results(iir_id ASC);

CREATE TABLE iir_drafts (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    data JSON NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_iir_drafts_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE UNIQUE INDEX idx_iir_drafts_user_unique ON iir_drafts(user_id ASC);
CREATE INDEX idx_iir_drafts_user_id ON iir_drafts(user_id ASC);
