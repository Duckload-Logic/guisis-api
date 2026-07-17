-- Rename programs table back to courses
RENAME TABLE programs TO courses;

-- Rename program_name column back to course_name
ALTER TABLE courses RENAME COLUMN program_name TO course_name;

-- Rename indexes in courses
ALTER TABLE courses RENAME INDEX unique_idx_program_code TO unique_idx_course_code;
ALTER TABLE courses RENAME INDEX unique_idx_program_name TO unique_idx_course_name;

-- Rename program_id column back to course_id
ALTER TABLE student_personal_info RENAME COLUMN program_id TO course_id;

-- Rename index on student_personal_info
ALTER TABLE student_personal_info RENAME INDEX idx_student_personal_info_program_id TO idx_student_personal_info_course_id;

-- Rename columns back in student_cors
ALTER TABLE student_cors RENAME COLUMN program_desc TO course_desc;
ALTER TABLE student_cors RENAME COLUMN program_code TO course_code;

-- Recreate original views referencing courses
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
    COALESCE(pf.file_url, '') AS two_by_two_photo_data_url,
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
JOIN iir_records iir ON iir.id = spi.iir_id
LEFT JOIN profile_pictures pp ON pp.user_id = iir.user_id
LEFT JOIN files pf ON pf.id = pp.file_id
LEFT JOIN civil_status_types cst ON spi.civil_status_id = cst.id
LEFT JOIN religions rel ON spi.religion_id = rel.id
LEFT JOIN courses c ON spi.course_id = c.id
LEFT JOIN student_statuses ss ON spi.status_id = ss.id
LEFT JOIN emergency_contacts ec ON spi.iir_id = ec.iir_id
LEFT JOIN student_relationship_types ert ON ec.relationship_id = ert.id;

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
    COALESCE(pf.file_url, '') AS profile_picture,
    COALESCE(c.code, '') AS course_code,
    COALESCE(c.course_name, '') AS course_name,
    iir.created_at,
    iir.updated_at
FROM iir_records iir
JOIN users u ON iir.user_id = u.id
JOIN student_personal_info spi ON iir.id = spi.iir_id
LEFT JOIN profile_pictures pp ON pp.user_id = iir.user_id
LEFT JOIN files pf ON pf.id = pp.file_id
LEFT JOIN student_statuses ss ON spi.status_id = ss.id
LEFT JOIN courses c ON spi.course_id = c.id
WHERE iir.is_submitted = TRUE;

CREATE OR REPLACE VIEW v_student_current_cors AS
SELECT
    sc.file_id,
    sc.student_id,
    sc.student_number,
    sc.course_desc,
    sc.course_code,
    sc.year_level,
    sc.section,
    sc.campus,
    sc.year_start,
    sc.year_end,
    sc.term,
    sc.valid_from,
    sc.valid_until
FROM
    student_cors sc
JOIN
    academic_settings ac ON ac.id = 1
WHERE
    sc.year_start = ac.current_year_start
    AND sc.term = ac.current_term
    AND sc.valid_from IS NOT NULL
    AND sc.valid_until IS NOT NULL;
