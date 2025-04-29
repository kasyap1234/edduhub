BEGIN;

DROP INDEX IF EXISTS idx_lectures_college_id;
DROP INDEX IF EXISTS idx_lectures_course_id;
-- DROP INDEX IF EXISTS idx_lectures_qr_code_id;
DROP INDEX IF EXISTS idx_lectures_datetime;
DROP TABLE IF EXISTS lectures;

COMMIT;
