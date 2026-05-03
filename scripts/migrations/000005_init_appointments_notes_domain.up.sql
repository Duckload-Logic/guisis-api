-- ============================================================================
-- APPOINTMENTS & SIGNIFICANT NOTES DOMAIN
-- ============================================================================

CREATE TABLE appointments (
    id CHAR(36) NOT NULL PRIMARY KEY,
    iir_id CHAR(36) NULL DEFAULT NULL,
    time_slot_id INT NOT NULL,
    when_date DATE NOT NULL,
    reason TEXT DEFAULT NULL,
    admin_notes TEXT DEFAULT NULL,
    appointment_category_id INT NOT NULL,
    urgency_level ENUM('LOW', 'MEDIUM', 'HIGH', 'CRITICAL') NOT NULL DEFAULT 'MEDIUM',
    urgency_score FLOAT NOT NULL,
    status_id INT NOT NULL DEFAULT 1,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT appointments_ibfk_2 FOREIGN KEY (time_slot_id) REFERENCES time_slots(id),
    CONSTRAINT appointments_ibfk_3 FOREIGN KEY (status_id) REFERENCES statuses(id),
    CONSTRAINT appointments_ibfk_4 FOREIGN KEY (appointment_category_id) REFERENCES appointment_categories(id),
    CONSTRAINT fk_appointments_iir FOREIGN KEY (iir_id) REFERENCES iir_records(id) ON DELETE SET NULL,
    UNIQUE KEY unique_appointment (when_date, time_slot_id)
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE UNIQUE INDEX unique_idx_appointment ON appointments(when_date ASC, time_slot_id ASC);
CREATE INDEX idx_appointments_time_slot_id ON appointments(time_slot_id ASC);
CREATE INDEX idx_appointments_appointment_category_id ON appointments(appointment_category_id ASC);
CREATE INDEX idx_appointments_status_id ON appointments(status_id ASC);
CREATE INDEX idx_appointments_when_date ON appointments(when_date ASC);
CREATE INDEX idx_appointments_iir_id ON appointments(iir_id ASC);
