BEGIN;

DROP INDEX IF EXISTS idx_courses_instructor_id;
DROP INDEX IF EXISTS idx_courses_name;
DROP TABLE IF EXISTS courses;

COMMIT;
