-- Drop unique key and unique index on appointments table
-- to allow rescheduling/new appointments on slots
-- where previous appointments were cancelled.
ALTER TABLE appointments DROP INDEX unique_appointment;
ALTER TABLE appointments DROP INDEX unique_idx_appointment;
