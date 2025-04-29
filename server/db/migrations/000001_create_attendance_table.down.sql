BEGIN;

DROP INDEX IF EXISTS idx_attendance_student_id;
DROP INDEX IF EXISTS idx_attendance_course_id;
DROP INDEX IF EXISTS idx_attendance_college_id;
DROP INDEX IF EXISTS idx_attendance_lecture_id;
DROP INDEX IF EXISTS idx_attendance_date;
DROP INDEX IF EXISTS idx_attendance_unique_key;
DROP TABLE IF EXISTS attendance;

COMMIT;
