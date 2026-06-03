CREATE OR REPLACE VIEW v_student_personal_info AS
SELECT
    spi.id,
    spi.iir_id,
    spi.student_number,
    spi.gender_id,
    COALESCE(g.gender_name, '') AS gender_name,
    spi.civil_status_id,
    COALESCE(cst.status_name, '') AS civil_status_name,
    spi.religion_id,
    COALESCE(rel.religion_name, '') AS religion_name,
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
LEFT JOIN genders g ON spi.gender_id = g.id
LEFT JOIN civil_status_types cst ON spi.civil_status_id = cst.id
LEFT JOIN religions rel ON spi.religion_id = rel.id
LEFT JOIN courses c ON spi.course_id = c.id
LEFT JOIN student_statuses ss ON spi.status_id = ss.id
LEFT JOIN emergency_contacts ec ON spi.iir_id = ec.iir_id
LEFT JOIN student_relationship_types ert ON ec.relationship_id = ert.id;

ALTER TABLE student_personal_info
DROP COLUMN other_religion_text;
