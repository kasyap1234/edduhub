-- Correct content for: 000004_create_students_table.down.sql
BEGIN;

-- Drop indexes associated with the students table
DROP INDEX IF EXISTS idx_students_user_id;
DROP INDEX IF EXISTS idx_students_college_id;
DROP INDEX IF EXISTS idx_students_kratos_identity_id;
DROP INDEX IF EXISTS idx_students_roll_no;

-- Drop the students table itself
DROP TABLE IF EXISTS students;

COMMIT;
