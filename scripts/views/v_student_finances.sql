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
