-- ============================================================================
-- AUTH & ROLES DOMAIN
-- ============================================================================

CREATE TABLE roles (
    id INT NOT NULL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE TABLE users (
    id CHAR(36) NOT NULL PRIMARY KEY,
    email VARCHAR(100) NOT NULL,
    first_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(100) DEFAULT NULL,
    last_name VARCHAR(100) NOT NULL,
    suffix_name VARCHAR(50) DEFAULT NULL,
    password_hash VARCHAR(255) DEFAULT NULL,
    auth_type VARCHAR(20) NOT NULL DEFAULT 'native',
    is_active TINYINT(1) DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE TABLE user_roles (
    user_id CHAR(36) NOT NULL,
    role_id INT NOT NULL,
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    assigned_by CHAR(36) DEFAULT NULL,
    reason TEXT DEFAULT NULL,
    reference_id VARCHAR(100) DEFAULT NULL,
    PRIMARY KEY (user_id, role_id),
    CONSTRAINT fk_user_roles_user FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_roles_role FOREIGN KEY (role_id) REFERENCES roles(id) ON DELETE CASCADE,
    CONSTRAINT fk_user_roles_admin FOREIGN KEY (assigned_by) REFERENCES users(id) ON DELETE SET NULL
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE UNIQUE INDEX idx_users_email_auth_type ON users(email ASC, auth_type ASC);