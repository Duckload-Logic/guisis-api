CREATE OR REPLACE VIEW v_student_financial_supports AS
SELECT
    sfs.sf_id,
    sst.id,
    sst.support_type_name
FROM
    student_financial_supports sfs
JOIN
    student_support_types sst ON sfs.support_type_id = sst.id;
