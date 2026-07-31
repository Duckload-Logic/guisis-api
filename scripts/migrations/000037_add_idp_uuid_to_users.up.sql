ALTER TABLE users ADD COLUMN idp_uuid VARCHAR(255) DEFAULT NULL;
CREATE UNIQUE INDEX idx_users_idp_uuid ON users(idp_uuid);
