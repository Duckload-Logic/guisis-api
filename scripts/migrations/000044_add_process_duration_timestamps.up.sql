-- Add process duration timestamps to appointments and admission slips

ALTER TABLE appointments
ADD COLUMN started_at TIMESTAMP NULL DEFAULT NULL AFTER status_id,
ADD COLUMN completed_at TIMESTAMP NULL DEFAULT NULL AFTER started_at;

ALTER TABLE admission_slips
ADD COLUMN started_at TIMESTAMP NULL DEFAULT NULL AFTER status_id,
ADD COLUMN completed_at TIMESTAMP NULL DEFAULT NULL AFTER started_at;
