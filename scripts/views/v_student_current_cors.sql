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
