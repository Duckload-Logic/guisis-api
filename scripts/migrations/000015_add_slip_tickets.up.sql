CREATE TABLE admission_tickets (
    id CHAR(36) NOT NULL PRIMARY KEY,
    admission_slip_id CHAR(36) NOT NULL,
    ticket_code VARCHAR(20) NOT NULL UNIQUE,
    is_verified BOOLEAN DEFAULT FALSE,
    verified_at TIMESTAMP NULL,
    verified_by CHAR(36) NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_tickets_slip FOREIGN KEY (admission_slip_id) 
        REFERENCES admission_slips(id) ON DELETE CASCADE,
    CONSTRAINT fk_tickets_verified_by FOREIGN KEY (verified_by) 
        REFERENCES users(id) ON DELETE SET NULL
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE INDEX idx_tickets_slip_id ON admission_tickets(admission_slip_id);
CREATE INDEX idx_tickets_code ON admission_tickets(ticket_code);
