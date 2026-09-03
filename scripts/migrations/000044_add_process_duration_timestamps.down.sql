-- Drop process duration timestamps from appointments and admission slips

ALTER TABLE appointments
DROP COLUMN completed_at,
DROP COLUMN started_at;

ALTER TABLE admission_slips
DROP COLUMN completed_at,
DROP COLUMN started_at;
