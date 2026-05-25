-- ============================================================================
-- LOOKUP & REFERENCE DOMAIN
-- ============================================================================

CREATE TABLE regions (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL UNIQUE,
    code VARCHAR(10) NOT NULL UNIQUE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE UNIQUE INDEX unique_idx_region_name ON regions(name ASC);
CREATE UNIQUE INDEX unique_idx_region_code ON regions(code ASC);

CREATE TABLE provinces (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(10) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    region_code VARCHAR(10) NOT NULL,
    CONSTRAINT provinces_ibfk_1 FOREIGN KEY (region_code) REFERENCES regions(code)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE UNIQUE INDEX unique_idx_progince_code ON provinces(code ASC);
CREATE INDEX idx_province_region_code ON provinces(region_code ASC);

CREATE TABLE cities (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(10) NOT NULL UNIQUE,
    name VARCHAR(100) DEFAULT NULL,
    type VARCHAR(20) DEFAULT NULL,
    zip_code VARCHAR(10) DEFAULT NULL,
    district VARCHAR(50) DEFAULT NULL,
    province_code VARCHAR(10) DEFAULT NULL,
    region_code VARCHAR(10) DEFAULT NULL,
    CONSTRAINT cities_ibfk_1 FOREIGN KEY (province_code) REFERENCES provinces(code),
    CONSTRAINT cities_ibfk_2 FOREIGN KEY (region_code) REFERENCES regions(code)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE UNIQUE INDEX unique_idx_city_code ON cities(code ASC);
CREATE UNIQUE INDEX unique_idx_province_city ON cities(province_code ASC, name ASC);
CREATE INDEX idx_cities_region_code ON cities(region_code ASC);
CREATE INDEX idx_cities_province_code ON cities(province_code ASC);

CREATE TABLE barangays (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(10) NOT NULL UNIQUE,
    name VARCHAR(100) DEFAULT NULL,
    city_code VARCHAR(10) DEFAULT NULL,
    CONSTRAINT barangays_ibfk_1 FOREIGN KEY (city_code) REFERENCES cities(code)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE UNIQUE INDEX unique_idx_barangay_code ON barangays(code ASC);
CREATE UNIQUE INDEX unique_idx_city_barangay ON barangays(city_code ASC, name ASC);
CREATE INDEX idx_barangays_city_code ON barangays(city_code ASC);

CREATE TABLE addresses (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    region_code VARCHAR(10) DEFAULT NULL,
    province_code VARCHAR(10) DEFAULT NULL,
    city_code VARCHAR(10) DEFAULT NULL,
    barangay_code VARCHAR(10) NOT NULL,
    street_detail VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT addresses_ibfk_1 FOREIGN KEY (region_code) REFERENCES regions(code),
    CONSTRAINT addresses_ibfk_2 FOREIGN KEY (city_code) REFERENCES cities(code),
    CONSTRAINT addresses_ibfk_3 FOREIGN KEY (province_code) REFERENCES provinces(code),
    CONSTRAINT addresses_ibfk_4 FOREIGN KEY (barangay_code) REFERENCES barangays(code)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE INDEX idx_addresses_region_id ON addresses(region_code ASC);
CREATE INDEX idx_addresses_province_id ON addresses(province_code ASC);
CREATE INDEX idx_addresses_city_id ON addresses(city_code ASC);
CREATE INDEX idx_addresses_barangay_id ON addresses(barangay_code ASC);

CREATE TABLE activity_options (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(100) NOT NULL,
    category ENUM('academic', 'extra_curricular', 'both') NOT NULL,
    is_active TINYINT(1) DEFAULT 1
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE TABLE admission_slip_categories (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE UNIQUE INDEX unique_idx_admission_slip_category_name ON admission_slip_categories(name ASC);

CREATE TABLE statuses (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    status_type ENUM('appointment', 'slip', 'both') NOT NULL DEFAULT 'both'
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE UNIQUE INDEX unique_idx_status_name ON statuses(name ASC);

CREATE TABLE appointment_categories (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE UNIQUE INDEX unique_idx_appointment_category_name ON appointment_categories(name ASC);

CREATE TABLE time_slots (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    time TIME NOT NULL UNIQUE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE UNIQUE INDEX unique_idx_time_slot_time ON time_slots(time ASC);

CREATE TABLE civil_status_types (
    id INT NOT NULL PRIMARY KEY,
    status_name VARCHAR(50) NOT NULL UNIQUE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE UNIQUE INDEX unique_idx_civil_status_name ON civil_status_types(status_name ASC);

CREATE TABLE genders (
    id INT NOT NULL PRIMARY KEY,
    gender_name VARCHAR(50) NOT NULL UNIQUE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE UNIQUE INDEX unique_idx_gender_name ON genders(gender_name ASC);

CREATE TABLE religions (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    religion_name VARCHAR(100) NOT NULL UNIQUE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE UNIQUE INDEX unique_idx_religion_name ON religions(religion_name ASC);

CREATE TABLE educational_levels (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    level_name VARCHAR(50) NOT NULL UNIQUE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE UNIQUE INDEX unique_idx_level_name ON educational_levels(level_name ASC);

CREATE TABLE student_relationship_types (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    relationship_name VARCHAR(100) NOT NULL UNIQUE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE UNIQUE INDEX unique_idx_relationship_name ON student_relationship_types(relationship_name ASC);

CREATE TABLE enrollment_reasons (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    reason_text VARCHAR(100) NOT NULL UNIQUE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE UNIQUE INDEX unique_idx_reason_text ON enrollment_reasons(reason_text ASC);

CREATE TABLE nature_of_residence_types (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    residence_type_name VARCHAR(100) NOT NULL UNIQUE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE UNIQUE INDEX unique_idx_residence_type_name ON nature_of_residence_types(residence_type_name ASC);

CREATE TABLE parental_status_types (
    id INT NOT NULL PRIMARY KEY,
    status_name VARCHAR(50) NOT NULL UNIQUE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE UNIQUE INDEX unique_idx_parental_status_name ON parental_status_types(status_name ASC);

CREATE TABLE courses (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    code VARCHAR(20) NOT NULL UNIQUE,
    course_name VARCHAR(100) NOT NULL UNIQUE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE UNIQUE INDEX unique_idx_course_code ON courses(code ASC);
CREATE UNIQUE INDEX unique_idx_course_name ON courses(course_name ASC);

CREATE TABLE income_ranges (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    range_text VARCHAR(100) NOT NULL
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE TABLE sibling_support_types (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE TABLE student_support_types (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    support_type_name VARCHAR(100) NOT NULL UNIQUE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE UNIQUE INDEX unique_idx_student_support_type_name ON student_support_types(support_type_name ASC);
