CREATE TABLE profile_pictures (
    file_id CHAR(36) PRIMARY KEY,
    user_id CHAR(36) NOT NULL,
    FOREIGN KEY (file_id) REFERENCES files(id) ON DELETE CASCADE,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    UNIQUE (user_id)
);

CREATE INDEX idx_profile_pictures_user_id ON profile_pictures(user_id);