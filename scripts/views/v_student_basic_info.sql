CREATE OR REPLACE VIEW v_student_basic_info AS
SELECT
    iir.id AS iir_id,
    u.id AS user_id,
    u.email,
    u.first_name,
    u.middle_name,
    u.last_name,
    u.suffix_name
FROM
    users u
JOIN
    iir_records iir ON u.id = iir.user_id;
