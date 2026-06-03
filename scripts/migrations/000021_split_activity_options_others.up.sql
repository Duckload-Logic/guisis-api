UPDATE activity_options SET is_active = 0 WHERE name = 'Others' AND category = 'both';

INSERT INTO activity_options (name, category, is_active) VALUES
    ('Others', 'academic', 1),
    ('Others', 'extra_curricular', 1);
