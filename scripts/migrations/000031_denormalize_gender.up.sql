-- Add gender ENUM column
ALTER TABLE student_personal_info
ADD COLUMN gender ENUM('Male', 'Female') NOT NULL DEFAULT 'Male';

-- Populate gender column from genders table
UPDATE student_personal_info
INNER JOIN genders ON student_personal_info.gender_id = genders.id
SET student_personal_info.gender = genders.gender_name;

-- Drop foreign key and index
ALTER TABLE student_personal_info
DROP FOREIGN KEY student_personal_info_ibfk_5;

DROP INDEX idx_student_personal_info_gender_id ON student_personal_info;

-- Drop gender_id column
ALTER TABLE student_personal_info
DROP COLUMN gender_id;

-- Drop genders table
DROP TABLE genders;

-- Recreate views to map the new gender column
CREATE OR REPLACE VIEW v_student_profiles AS
SELECT
    iir.id AS iir_id,
    iir.user_id,
    u.first_name,
    u.middle_name,
    u.last_name,
    u.suffix_name,
    u.email,
    spi.student_number,
    CASE WHEN spi.gender = 'Male' THEN 1 ELSE 2 END AS gender_id,
    spi.course_id,
    spi.section,
    spi.year_level,
    spi.status_id,
    COALESCE(ss.status_name, '') AS status_name,
    spi.gender AS gender_name,
    COALESCE(c.code, '') AS course_code,
    COALESCE(c.course_name, '') AS course_name,
    iir.created_at,
    iir.updated_at
FROM iir_records iir
JOIN users u ON iir.user_id = u.id
JOIN student_personal_info spi ON iir.id = spi.iir_id
LEFT JOIN student_statuses ss ON spi.status_id = ss.id
LEFT JOIN courses c ON spi.course_id = c.id
WHERE iir.is_submitted = TRUE;

CREATE OR REPLACE VIEW v_student_personal_info AS
SELECT
    spi.id,
    spi.iir_id,
    spi.student_number,
    CASE WHEN spi.gender = 'Male' THEN 1 ELSE 2 END AS gender_id,
    spi.gender AS gender_name,
    spi.civil_status_id,
    COALESCE(cst.status_name, '') AS civil_status_name,
    spi.religion_id,
    COALESCE(rel.religion_name, '') AS religion_name,
    spi.other_religion_text,
    spi.height_m,
    spi.weight_kg,
    spi.complexion,
    spi.high_school_gwa,
    spi.course_id,
    COALESCE(c.code, '') AS course_code,
    COALESCE(c.course_name, '') AS course_name,
    spi.year_level,
    spi.section,
    spi.place_of_birth,
    spi.date_of_birth,
    spi.is_employed,
    spi.employer_name,
    spi.employer_address,
    spi.mobile_number,
    spi.telephone_number,
    spi.employer_contact_number,
    spi.status_id,
    COALESCE(ss.status_name, '') AS status_name,
    spi.graduation_year,
    COALESCE(ec.id, 0) AS emergency_id,
    COALESCE(ec.first_name, '') AS emergency_first_name,
    COALESCE(ec.middle_name, '') AS emergency_middle_name,
    COALESCE(ec.last_name, '') AS emergency_last_name,
    COALESCE(ec.contact_number, '') AS emergency_contact_number,
    COALESCE(ec.relationship_id, 0) AS emergency_relationship_id,
    COALESCE(ert.relationship_name, '') AS emergency_relationship_name,
    COALESCE(ec.address_id, 0) AS emergency_address_id
FROM student_personal_info spi
LEFT JOIN civil_status_types cst ON spi.civil_status_id = cst.id
LEFT JOIN religions rel ON spi.religion_id = rel.id
LEFT JOIN courses c ON spi.course_id = c.id
LEFT JOIN student_statuses ss ON spi.status_id = ss.id
LEFT JOIN emergency_contacts ec ON spi.iir_id = ec.iir_id
LEFT JOIN student_relationship_types ert ON ec.relationship_id = ert.id;
