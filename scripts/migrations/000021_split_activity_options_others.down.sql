DELETE FROM activity_options WHERE name = 'Others' AND category IN ('academic', 'extra_curricular');

UPDATE activity_options SET is_active = 1 WHERE name = 'Others' AND category = 'both';
