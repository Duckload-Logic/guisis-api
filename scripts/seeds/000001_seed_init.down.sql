SET FOREIGN_KEY_CHECKS = 0;

TRUNCATE TABLE appointment_types;
TRUNCATE TABLE health_remark_types;
TRUNCATE TABLE financial_support_types;
TRUNCATE TABLE parental_status_types;
TRUNCATE TABLE religion_types;
TRUNCATE TABLE civil_status_types;
TRUNCATE TABLE educational_levels;
TRUNCATE TABLE address_types;
TRUNCATE TABLE guardian_types;
TRUNCATE TABLE roles;

SET FOREIGN_KEY_CHECKS = 1;

DROP TABLE seed_migrations;