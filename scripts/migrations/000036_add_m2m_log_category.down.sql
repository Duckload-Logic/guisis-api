UPDATE system_logs
SET category = 'SYSTEM'
WHERE category = 'M2M';

ALTER TABLE system_logs MODIFY COLUMN category ENUM('SECURITY', 'SYSTEM', 'AUDIT', 'CONSENT') NOT NULL;
