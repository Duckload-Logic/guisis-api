-- ============================================================================
-- ADD STUDENT STATUS AND LOCKING
-- ============================================================================

CREATE TABLE student_statuses (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    status_name VARCHAR(50) NOT NULL UNIQUE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

INSERT INTO student_statuses (status_name) VALUES 
('Active'), 
('Graduated'), 
('On Leave'), 
('Archived'), 
('Withdrawn');

ALTER TABLE student_personal_info 
ADD COLUMN status_id INT NOT NULL DEFAULT 1,
ADD COLUMN graduation_year INT DEFAULT NULL;

ALTER TABLE student_personal_info
ADD CONSTRAINT fk_student_status FOREIGN KEY (status_id) REFERENCES student_statuses(id);

CREATE INDEX idx_student_personal_info_status_id ON student_personal_info(status_id ASC);
