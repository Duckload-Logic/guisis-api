-- ============================================================================
-- ADMISSION SLIPS DOMAIN
-- ============================================================================

CREATE TABLE admission_slips (
    id CHAR(36) NOT NULL PRIMARY KEY,
    iir_id CHAR(36) NOT NULL,
    category_id INT NOT NULL,
    reason TEXT NOT NULL,
    date_of_absence DATE NOT NULL,
    date_needed DATE NOT NULL,
    status_id INT NOT NULL DEFAULT 1,
    admin_notes TEXT DEFAULT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    CONSTRAINT admission_slips_ibfk_2 FOREIGN KEY (status_id) REFERENCES statuses(id),
    CONSTRAINT admission_slips_ibfk_3 FOREIGN KEY (category_id) REFERENCES admission_slip_categories(id),
    CONSTRAINT fk_admission_slips_iir FOREIGN KEY (iir_id) REFERENCES iir_records(id) ON DELETE CASCADE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE INDEX idx_admission_slips_status_id ON admission_slips(status_id ASC);
CREATE INDEX idx_admission_slips_category_id ON admission_slips(category_id ASC);
CREATE INDEX idx_admission_slips_iir_id ON admission_slips(iir_id ASC);

CREATE TABLE slip_attachments (
    file_id CHAR(36) NOT NULL PRIMARY KEY,
    admission_slip_id CHAR(36) NOT NULL,
    attachment_type ENUM('MEDICAL', 'EXCUSE LETTER', 'PARENT VALID ID', 'OTHER') NOT NULL,
    FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE,
    FOREIGN KEY (admission_slip_id) REFERENCES admission_slips(id) ON DELETE CASCADE
) ENGINE = InnoDB DEFAULT CHARACTER SET = utf8mb4 COLLATE = utf8mb4_0900_ai_ci;

CREATE INDEX idx_slip_attachments_admission_slip_id ON slip_attachments(admission_slip_id ASC);
