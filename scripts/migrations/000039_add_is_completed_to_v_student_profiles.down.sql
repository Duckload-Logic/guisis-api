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
    spi.program_id,
    spi.section,
    spi.year_level,
    spi.status_id,
    COALESCE(ss.status_name, '') AS status_name,
    spi.gender AS gender_name,
    COALESCE(pf.file_url, '') AS profile_picture,
    COALESCE(p.code, '') AS program_code,
    COALESCE(p.program_name, '') AS program_name,
    iir.created_at,
    iir.updated_at
FROM iir_records iir
JOIN users u ON iir.user_id = u.id
JOIN student_personal_info spi ON iir.id = spi.iir_id
LEFT JOIN profile_pictures pp ON pp.user_id = iir.user_id
LEFT JOIN files pf ON pf.id = pp.file_id
LEFT JOIN student_statuses ss ON spi.status_id = ss.id
LEFT JOIN programs p ON spi.program_id = p.id
WHERE iir.is_submitted = TRUE;
