ALTER TABLE academic_settings
ADD COLUMN allow_expedited_iir TINYINT(1) NOT NULL DEFAULT 0;

ALTER TABLE iir_records
ADD COLUMN is_completed TINYINT(1) NOT NULL DEFAULT 1;
