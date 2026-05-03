-- ============================================================================
-- SYSTEM DOMAIN
-- ============================================================================

CREATE TABLE m2m_clients (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    client_name VARCHAR(100) NOT NULL,
    client_id VARCHAR(36) NOT NULL UNIQUE,
    client_secret_hash VARCHAR(64) NOT NULL,
    client_description VARCHAR(255) NOT NULL,
    scopes JSON DEFAULT NULL,
    is_active TINYINT(1) NOT NULL DEFAULT 1,
    is_verified TINYINT(1) NOT NULL DEFAULT 0,
    last_used_at TIMESTAMP NULL DEFAULT NULL,
    expires_at TIMESTAMP NULL DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE UNIQUE INDEX unique_idx_m2m_clients_client_id ON m2m_clients(client_id ASC);
CREATE INDEX idx_m2m_clients_user_id ON m2m_clients(user_id ASC);
CREATE INDEX idx_m2m_clients_is_active ON m2m_clients(is_active ASC);

CREATE TABLE notifications (
    id CHAR(36) NOT NULL PRIMARY KEY,
    receiver_id CHAR(36) NOT NULL,
    actor_id CHAR(36) NULL DEFAULT NULL,
    -- TargetID is now a "Soft Link" (No Foreign Key Constraint)
    target_id CHAR(36) NULL DEFAULT NULL,
    -- TargetType tells the code if this is an 'Appointment', 'Slip', etc.
    target_type VARCHAR(50) NULL DEFAULT NULL,
    title VARCHAR(255) NOT NULL,
    message TEXT NOT NULL,
    type ENUM('Appointment', 'Slip', 'Guidance', 'System', 'General')
        DEFAULT 'System',
    is_read TINYINT(1) DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,

    -- Keep constraints for the people involved
    CONSTRAINT fk_notifications_receiver
        FOREIGN KEY (receiver_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_notifications_actor
        FOREIGN KEY (actor_id) REFERENCES users(id) ON DELETE SET NULL
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE INDEX idx_notifications_receiver_id ON notifications(receiver_id ASC);
CREATE INDEX idx_notifications_actor_id ON notifications(actor_id ASC);

CREATE TABLE system_logs (
    id INT NOT NULL AUTO_INCREMENT PRIMARY KEY,
    level ENUM('INFO', 'WARNING', 'ERROR', 'CRITICAL') NOT NULL,
    category ENUM('SECURITY', 'SYSTEM', 'AUDIT', 'CONSENT') NOT NULL,
    action VARCHAR(100) NOT NULL,
    message TEXT NOT NULL,

    -- IDs as Strings (No Foreign Keys for maximum persistence)
    user_id CHAR(36) NULL DEFAULT NULL,
    target_id CHAR(36) NULL DEFAULT NULL,
    target_type VARCHAR(50) NULL DEFAULT NULL, -- 'User', 'Appt', 'Slip'

    -- Denormalized data (Snapshot of the moment it happened)
    user_email VARCHAR(100) NULL DEFAULT NULL,
    target_email VARCHAR(100) NULL DEFAULT NULL,

    ip_address VARCHAR(45) NULL DEFAULT NULL,
    user_agent VARCHAR(255) NULL DEFAULT NULL,
    metadata JSON DEFAULT NULL,
    trace_id CHAR(36) NULL DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE INDEX idx_system_logs_action ON system_logs(action ASC);
CREATE INDEX idx_system_logs_user_email ON system_logs(user_email ASC);
CREATE INDEX idx_system_logs_target_email ON system_logs(target_email ASC);
CREATE INDEX idx_system_logs_created_at ON system_logs(created_at ASC);
CREATE INDEX idx_system_logs_category_created ON system_logs(created_at ASC);
CREATE INDEX idx_system_logs_user_id ON system_logs(user_id ASC);
