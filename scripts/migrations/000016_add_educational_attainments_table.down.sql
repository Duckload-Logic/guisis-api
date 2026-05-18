ALTER TABLE related_persons DROP FOREIGN KEY fk_related_persons_attainment;
ALTER TABLE related_persons DROP COLUMN educational_attainment_id;
ALTER TABLE related_persons ADD COLUMN educational_level VARCHAR(100) NOT NULL AFTER suffix_name;
DROP TABLE educational_attainments;
