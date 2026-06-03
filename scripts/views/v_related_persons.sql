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
FROM
    student_related_persons srp
JOIN
    related_persons rp ON srp.related_person_id = rp.id
LEFT JOIN
    educational_attainments ea ON rp.educational_attainment_id = ea.id
LEFT JOIN
    student_relationship_types ert ON srp.relationship_id = ert.id;
