-- Restore unique key and unique index on appointments table
ALTER TABLE appointments ADD CONSTRAINT unique_appointment UNIQUE (when_date, time_slot_id);
CREATE UNIQUE INDEX unique_idx_appointment ON appointments(when_date ASC, time_slot_id ASC);
