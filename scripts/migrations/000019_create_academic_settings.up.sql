CREATE TABLE academic_settings (
    id INT NOT NULL PRIMARY KEY DEFAULT 1,
    current_year_start INT NOT NULL,
    current_year_end INT NOT NULL,
    current_term INT NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
        ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT chk_academic_settings_id CHECK (id = 1)
) ENGINE = InnoDB
  DEFAULT CHARACTER SET = utf8mb4
  COLLATE = utf8mb4_0900_ai_ci;

INSERT INTO academic_settings
    (id, current_year_start, current_year_end, current_term)
VALUES (1, 2025, 2026, 1);
