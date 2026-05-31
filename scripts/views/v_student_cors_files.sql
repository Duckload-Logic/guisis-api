CREATE OR REPLACE VIEW v_student_cors_files AS
SELECT
    sc.student_id,
    f.file_url
FROM
    student_cors sc
JOIN
    files f ON f.id = sc.file_id;
