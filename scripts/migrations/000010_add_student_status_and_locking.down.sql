-- ============================================================================
-- DOWN: REMOVE STUDENT STATUS AND LOCKING
-- ============================================================================

ALTER TABLE student_personal_info DROP FOREIGN KEY fk_student_status;
DROP INDEX idx_student_personal_info_status_id ON student_personal_info;
ALTER TABLE student_personal_info DROP COLUMN status_id;
ALTER TABLE student_personal_info DROP COLUMN graduation_year;

DROP TABLE student_statuses;
