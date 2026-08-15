ALTER TABLE student_activities ADD COLUMN role VARCHAR(255) DEFAULT 'Member';

UPDATE student_activities sa
JOIN (
  SELECT student_activity_id, GROUP_CONCAT(role ORDER BY role SEPARATOR ', ')
    AS roles
  FROM student_activity_roles
  GROUP BY student_activity_id
) r ON sa.id = r.student_activity_id
SET sa.role = r.roles;

DROP TABLE IF EXISTS student_activity_roles;
