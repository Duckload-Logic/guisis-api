CREATE TABLE student_activity_roles (
  id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
  student_activity_id INT NOT NULL,
  role ENUM('Officer', 'Member', 'Other') NOT NULL,
  created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
  CONSTRAINT fk_student_activity_roles_activity
    FOREIGN KEY (student_activity_id)
    REFERENCES student_activities(id) ON DELETE CASCADE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

INSERT INTO student_activity_roles (student_activity_id, role)
SELECT id, 'Member' FROM student_activities WHERE role LIKE '%Member%';

INSERT INTO student_activity_roles (student_activity_id, role)
SELECT id, 'Officer' FROM student_activities WHERE role LIKE '%Officer%';

INSERT INTO student_activity_roles (student_activity_id, role)
SELECT id, 'Other' FROM student_activities WHERE role LIKE '%Other%';

ALTER TABLE student_activities DROP COLUMN role;
