-- Add Student Assistant Role
INSERT INTO roles (id, name)
VALUES (5, 'Student Assistant')
ON DUPLICATE KEY UPDATE name = VALUES(name);
