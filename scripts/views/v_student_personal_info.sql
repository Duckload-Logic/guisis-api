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
    spi.program_id,
    COALESCE(p.code, '') AS program_code,
    COALESCE(p.program_name, '') AS program_name,
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
LEFT JOIN programs p ON spi.program_id = p.id
LEFT JOIN student_statuses ss ON spi.status_id = ss.id
LEFT JOIN emergency_contacts ec ON spi.iir_id = ec.iir_id
LEFT JOIN student_relationship_types ert ON ec.relationship_id = ert.id;
