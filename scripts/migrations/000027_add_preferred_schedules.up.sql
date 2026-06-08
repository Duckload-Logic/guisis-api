ALTER TABLE appointments
ADD COLUMN preferred_date_1 DATE DEFAULT NULL,
ADD COLUMN preferred_time_slot_id_1 INT DEFAULT NULL,
ADD COLUMN preferred_date_2 DATE DEFAULT NULL,
ADD COLUMN preferred_time_slot_id_2 INT DEFAULT NULL,
ADD COLUMN preferred_date_3 DATE DEFAULT NULL,
ADD COLUMN preferred_time_slot_id_3 INT DEFAULT NULL,
ADD CONSTRAINT fk_appointments_pref_time_slot_1 FOREIGN KEY (preferred_time_slot_id_1) REFERENCES time_slots(id) ON DELETE SET NULL,
ADD CONSTRAINT fk_appointments_pref_time_slot_2 FOREIGN KEY (preferred_time_slot_id_2) REFERENCES time_slots(id) ON DELETE SET NULL,
ADD CONSTRAINT fk_appointments_pref_time_slot_3 FOREIGN KEY (preferred_time_slot_id_3) REFERENCES time_slots(id) ON DELETE SET NULL;
