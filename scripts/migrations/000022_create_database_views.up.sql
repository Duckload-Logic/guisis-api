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

CREATE OR REPLACE VIEW v_student_finances AS
SELECT
    sf.id,
    sf.iir_id,
    sf.monthly_family_income_range_id AS income_range_id,
    COALESCE(ir.range_text, '') AS income_range_text,
    sf.other_income_details AS other_income,
    sf.weekly_allowance
FROM student_finances sf
LEFT JOIN income_ranges ir
    ON sf.monthly_family_income_range_id = ir.id;

CREATE OR REPLACE VIEW v_related_persons AS
SELECT
    rp.id,
    rp.last_name,
    rp.first_name,
    rp.middle_name,
    rp.suffix_name,
    rp.date_of_birth,
    COALESCE(rp.educational_attainment_id, 0) AS educational_attainment_id,
    COALESCE(ea.name, '') AS educational_attainment_name,
    COALESCE(rp.occupation, '') AS occupation,
    COALESCE(rp.employer_name, '') AS employer_name,
    COALESCE(rp.employer_address, '') AS employer_address,
    srp.iir_id,
    srp.relationship_id,
    COALESCE(ert.relationship_name, '') AS relationship_name,
    srp.is_parent,
    srp.is_guardian,
    srp.is_living
FROM student_related_persons srp
JOIN related_persons rp ON srp.related_person_id = rp.id
LEFT JOIN educational_attainments ea
    ON rp.educational_attainment_id = ea.id
LEFT JOIN student_relationship_types ert ON srp.relationship_id = ert.id;
