ALTER TABLE notifications ADD COLUMN is_touched TINYINT(1) DEFAULT 0 AFTER is_read;
