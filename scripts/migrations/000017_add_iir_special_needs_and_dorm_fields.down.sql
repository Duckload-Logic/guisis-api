ALTER TABLE student_personal_info
DROP COLUMN employer_contact_number,
DROP COLUMN living_in_dorm,
DROP COLUMN dorm_address,
DROP COLUMN landlord_name,
DROP COLUMN landlord_contact_number;

ALTER TABLE student_health_records
DROP COLUMN mental_emotional_has_problem,
DROP COLUMN mental_emotional_details;
