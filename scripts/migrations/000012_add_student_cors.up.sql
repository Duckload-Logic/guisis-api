CREATE TABLE student_cors (
    file_id CHAR(36) PRIMARY KEY,
    student_id CHAR(36) NOT NULL,
    valid_from DATE,
    valid_until DATE,
    FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE,
    FOREIGN KEY (student_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX idx_student_cors_student_id ON student_cors(student_id);