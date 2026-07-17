INSERT IGNORE INTO roles (id, `name`)
VALUES
    (1, 'STUDENT'),
    (2, 'ADMIN'),
    (3, 'SUPERADMIN'),
    (4, 'DEVELOPER');

INSERT INTO student_support_types (support_type_name) VALUES
    ('Parents'),
    ('Brother/Sister'),
    ('Spouse'),
    ('Scholarship'),
    ('Relatives'),
    ('Self-supporting/working student');

INSERT INTO income_ranges (id, range_text) VALUES
    (1, 'Below Php 5,000'),
    (2, 'Php 5,001 - Php 10,000'),
    (3, 'Php 10,001 - Php 15,000'),
    (4, 'Php 15,001 - Php 20,000'),
    (5, 'Php 20,001 - Php 30,000'),
    (6, 'Php 30,001 - Php 35,000'),
    (7, 'Php 35,001 - Php 40,000'),
    (8, 'Php 40,001 - Php 45,000'),
    (9, 'Php 45,001 - Php 50,000'),
    (10, 'Above Php 50,001');

INSERT INTO parental_status_types (id, status_name) VALUES
    (1, 'Married and staying together'),
    (2, 'Not Married but living together'),
    (3, 'Single Parent'),
    (4, 'Married but Separated'),
    (5, 'Other');

INSERT INTO educational_levels (level_name) VALUES
    ('Pre-Elementary'),
    ('Elementary'),
    ('Junior High School'),
    ('Senior High School'),
    ('Vocational'),
    ('College');

INSERT INTO programs (code, program_name) VALUES
    ('BSBA-HRM', 'Bachelor of Science in Business Administration - Human Resource Management'),
    ('BSBA-MM', 'Bachelor of Science in Business Administration - Marketing Management'),
    ('BSED-ENGLISH', 'Bachelor of Science in Education - English'),
    ('BSED-MATH', 'Bachelor of Science in Education - Mathematics'),
    ('BSECE', 'Bachelor of Science in Electronics and Communications Engineering'),
    ('BSIT', 'Bachelor of Science in Information Technology'),
    ('BSME', 'Bachelor of Science in Mechanical Engineering'),
    ('BSOA', 'Bachelor of Science in Office Administration'),
    ('BSPSYCH', 'Bachelor of Science in Psychology'),
    ('DIT', 'Diploma in Information Technology'),
    ('DOMT', 'Diploma in Office Management Technology');

INSERT INTO civil_status_types (id, status_name) VALUES
    (1, 'Single'),
    (2, 'Married'),
    (3, 'Widowed'),
    (4, 'Separated'),
    (5, 'Divorced');

INSERT INTO student_relationship_types (relationship_name) VALUES
    ('Father'),
    ('Mother'),
    ('Uncle'),
    ('Auntie'),
    ('Brother'),
    ('Sister'),
    ('Grandmother'),
    ('Grandfather'),
    ('Cousin'),
    ('Legal Guardian'),
    ('Friend'),
    ('Other relative');

INSERT INTO nature_of_residence_types (residence_type_name) VALUES
    ('Family home'),
    ("Relative's house"),
    ('Bed spacer'),
    ('House of married brother/sister'),
    ('Rented apartment/house'),
    ('Dormitory'),
    ('Shares apartment with friends/relatives');

INSERT INTO religions (religion_name) VALUES
    ('Christian/Born Again'),
    ('Roman Catholic'),
    ('Baptist'),
    ('Iglesia ni Cristo'),
    ('Islam/Muslim'),
    ('Protestant'),
    ('Members Church of God International'),
    ("Jehovah's Witness"),
    ('Seventh-Day Adventist'),
    ('Mormons/Latter-day Saints'),
    ('Apostolic'),
    ('United Church of Christ in the Philippines'),
    ('Church of Christ'),
    ('United Pentecostal'),
    ('Others'),
    ('Not Indicated');

INSERT INTO sibling_support_types (`name`) VALUES
    ('Family'),
    ('Your studies'),
    ('His/Her own family');

INSERT INTO activity_options (`name`, category) VALUES
    ('Math Club', 'academic'),
    ('Science Club', 'academic'),
    ('Debating Club', 'academic'),
    ("Quizzer's Club", 'academic'),
    ('Athletics', 'extra_curricular'),
    ('Dramatics', 'extra_curricular'),
    ('Religious Organizations', 'extra_curricular'),
    ('Chess Club', 'extra_curricular'),
    ('Glee Club', 'extra_curricular'),
    ('Scouting', 'extra_curricular'),
    ('Others', 'academic'),
    ('Others', 'extra_curricular');

INSERT INTO appointment_categories (id, name) VALUES
    (1, 'Academic'),
    (2, 'Financial'),
    (3, 'Personal'),
    (4, 'Career Guidance'),
    (5, 'Mental Health'),
    (6, 'Psychological Testing'),
    (7, 'Other');

INSERT INTO admission_slip_categories (id, name) VALUES
    (1, 'Medical'),
    (2, 'Personal'),
    (3, 'Family-related'),
    (4, 'Scholarship'),
    (5, 'Other');

INSERT INTO statuses (id, name, status_type) VALUES
    (1, 'Pending', 'both'),
    (2, 'Scheduled', 'appointment'),
    (3, 'Completed',  'appointment'),
    (4, 'Cancelled', 'appointment'),
    (5, 'Rejected', 'both'),
    (6, 'Rescheduled', 'appointment'),
    (7, 'No-show', 'appointment'),
    (8, 'Approved', 'slip'),
    (9, 'For Revision', 'slip');

INSERT INTO time_slots (id, time) VALUES
    (1, '08:00:00'),
    (2, '09:00:00'),
    (3, '10:00:00'),
    (4, '11:00:00'),
    (5, '13:00:00'),
    (6, '14:00:00'),
    (7, '15:00:00'),
    (8, '16:00:00');
