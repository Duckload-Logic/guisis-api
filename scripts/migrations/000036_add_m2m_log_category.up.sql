ALTER TABLE system_logs MODIFY COLUMN category ENUM('SECURITY', 'SYSTEM', 'AUDIT', 'CONSENT', 'M2M') NOT NULL;

UPDATE system_logs
SET category = 'M2M'
WHERE action LIKE 'M2M_%';
