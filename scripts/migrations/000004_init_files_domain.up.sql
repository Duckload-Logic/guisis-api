CREATE TABLE files (
    id CHAR(36) PRIMARY KEY,
    file_name VARCHAR(255) NOT NULL,
    file_url VARCHAR(255) NOT NULL,
    file_type VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    mime_type VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP NULL DEFAULT NULL
);

CREATE INDEX idx_files_id ON files(id);
CREATE INDEX idx_files_created_at ON files(created_at);
CREATE INDEX idx_files_deleted_at ON files(deleted_at);