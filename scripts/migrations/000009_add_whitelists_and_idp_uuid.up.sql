
CREATE TABLE whitelists (
    email VARCHAR(100) NOT NULL,
    role_id INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (email, role_id),
    CONSTRAINT whitelists_ibfk_1 FOREIGN KEY (role_id) REFERENCES roles(id)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;
