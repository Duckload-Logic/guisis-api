CREATE TABLE file_ocr_results (
    file_id CHAR(36) PRIMARY KEY,
    raw_text LONGTEXT,
    structured_data JSON,
    engine_v VARCHAR(50),
    confidence FLOAT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE
);
