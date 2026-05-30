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
    spi.gender_id,
    spi.course_id,
    spi.section,
    spi.year_level,
    spi.status_id,
    COALESCE(ss.status_name, '') AS status_name,
    COALESCE(g.gender_name, '') AS gender_name,
    COALESCE(c.code, '') AS course_code,
    COALESCE(c.course_name, '') AS course_name,
    iir.created_at,
    iir.updated_at
FROM iir_records iir
JOIN users u ON iir.user_id = u.id
JOIN student_personal_info spi ON iir.id = spi.iir_id
LEFT JOIN student_statuses ss ON spi.status_id = ss.id
LEFT JOIN genders g ON spi.gender_id = g.id
LEFT JOIN courses c ON spi.course_id = c.id
WHERE iir.is_submitted = TRUE;
