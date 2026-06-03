CREATE TABLE educational_attainments (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    name VARCHAR(255) NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

-- Seed initial data
INSERT INTO educational_attainments (name) VALUES
("Doctorate Degree"),
("Master's Degree"),
('College Undergraduate'),
('College Graduate'),
('Vocational'),
('High School Undergraduate'),
('High School Graduate'),
('Elementary Undergraduate'),
('Elementary Graduate'),
('Not Indicated');

-- Add column to related_persons and remove obsolete educational_level
ALTER TABLE related_persons ADD COLUMN educational_attainment_id INT;
ALTER TABLE related_persons ADD CONSTRAINT fk_related_persons_attainment FOREIGN KEY (educational_attainment_id) REFERENCES educational_attainments(id);
ALTER TABLE related_persons DROP COLUMN educational_level;
