BEGIN;

DROP INDEX IF EXISTS idx_enrollments_student_id;
DROP INDEX IF EXISTS idx_enrollments_course_id;
DROP INDEX IF EXISTS idx_enrollments_college_id;
DROP INDEX IF EXISTS idx_enrollments_student_course;
DROP TABLE IF EXISTS enrollments;

COMMIT;
