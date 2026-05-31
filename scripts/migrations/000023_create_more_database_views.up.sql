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

CREATE OR REPLACE VIEW v_student_financial_supports AS
SELECT
    sfs.sf_id,
    sst.id,
    sst.support_type_name
FROM
    student_financial_supports sfs
JOIN
    student_support_types sst ON sfs.support_type_id = sst.id;

CREATE OR REPLACE VIEW v_student_cors_files AS
SELECT
    sc.student_id,
    f.file_url
FROM
    student_cors sc
JOIN
    files f ON f.id = sc.file_id;

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
