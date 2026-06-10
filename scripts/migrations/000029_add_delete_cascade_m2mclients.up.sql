    ALTER TABLE
        m2m_clients
    ADD CONSTRAINT
        fk_m2m_user
    FOREIGN KEY
        (user_id)
    REFERENCES
        users(id)
    ON DELETE CASCADE;