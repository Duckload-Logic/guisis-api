ALTER TABLE appointments
DROP FOREIGN KEY fk_appointments_pref_time_slot_1,
DROP FOREIGN KEY fk_appointments_pref_time_slot_2,
DROP FOREIGN KEY fk_appointments_pref_time_slot_3,
DROP COLUMN preferred_date_1,
DROP COLUMN preferred_time_slot_id_1,
DROP COLUMN preferred_date_2,
DROP COLUMN preferred_time_slot_id_2,
DROP COLUMN preferred_date_3,
DROP COLUMN preferred_time_slot_id_3;
