CREATE TABLE significant_notes (
    id CHAR(36) NOT NULL PRIMARY KEY,
    iir_id CHAR(36) NULL DEFAULT NULL,
    appointment_id CHAR(36) NULL DEFAULT NULL,
    admission_slip_id CHAR(36) NULL DEFAULT NULL,
    note TEXT DEFAULT NULL,
    remarks TEXT DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT fk_sig_notes_iir FOREIGN KEY (iir_id) REFERENCES iir_records(id) ON DELETE SET NULL,
    CONSTRAINT fk_sig_notes_appointment FOREIGN KEY (appointment_id) REFERENCES appointments(id) ON DELETE CASCADE,
    CONSTRAINT fk_sig_notes_admission_slip FOREIGN KEY (admission_slip_id) REFERENCES admission_slips(id) ON DELETE CASCADE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE INDEX idx_significant_notes_iir_id ON significant_notes(iir_id ASC);
CREATE INDEX idx_significant_notes_appointment_id ON significant_notes(appointment_id ASC);
CREATE INDEX idx_significant_notes_admission_slip_id ON significant_notes(admission_slip_id ASC);
