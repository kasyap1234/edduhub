package repository

import (
	"context"
	"fmt" // Import fmt for better error wrapping
	"time"

	// Assuming models.Attendance uses time.Time
	"eduhub/server/internal/models" // Your models package

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan" // Import pgxscan
)

type AttendanceRepository interface {
	GetAttendanceByCourse(ctx context.Context, collegeID int, courseID int) ([]*models.Attendance, error)
	MarkAttendance(ctx context.Context, collegeID int, studentID int, courseID int, lectureID int) (bool, error)
	UpdateAttendance(ctx context.Context, collegeID int, studentID int, courseID int, lectureID int, status string) error
	GetAttendanceStudentInCourse(ctx context.Context, collegeID int, studentID int, courseID int) ([]*models.Attendance, error)
	GetAttendanceStudent(ctx context.Context, collegeID int, studentID int) ([]*models.Attendance, error)
	GetAttendanceByLecture(ctx context.Context, collegeID int, lectureID int, courseID int) ([]*models.Attendance, error)
	FreezeAttendance(ctx context.Context, collegeID int, studentID int) error
	// ProcessQRCode(ctx context.Context, collegeID int, studentID int, courseID int, lectureID int) (bool, error)
}

const attendanceTable = "attendance"

type attendanceRepository struct {
	DB *DB // Assuming DB struct is accessible here
}

// Assuming models.Attendance struct with db tags is defined in models package:
// type Attendance struct {
//     ID        int       `db:"id"`
//     StudentID int       `db:"student_id"`
//     CourseID  int       `db:"course_id"`
//     CollegeID int       `db:"college_id"`
//     Date      time.Time `db:"date"`
//     Status    string    `db:"status"`
//     ScannedAt time.Time `db:"scanned_at"`
//     LectureID int       `db:"lecture_id"`
// }

func NewAttendanceRepository(DB *DB) AttendanceRepository {
	return &attendanceRepository{
		DB: DB,
	}
}

func (a *attendanceRepository) GetAttendanceByCourse(
	ctx context.Context,
	collegeID int,
	courseID int,
) ([]*models.Attendance, error) {
	// Define the table name (assuming it's "attendance")
	const attendanceTable = "attendance"

	// Build the SELECT query
	query := a.DB.SQ.Select(
		"id", // Use database column names matching struct tags
		"student_id",
		"course_id",
		"college_id",
		"date",
		"status",
		"scanned_at",
		"lecture_id",
	).
		From(attendanceTable). // Specify the table
		Where(squirrel.Eq{     // Use WHERE to filter
			// Use database column names matching struct tags
			"college_id": collegeID,
			"course_id":  courseID,
		})

	sql, args, err := query.ToSql()
	if err != nil {
		// Use fmt.Errorf to wrap the original error for better debugging
		return nil, fmt.Errorf("GetAttendanceByCourse: failed to build query: %w", err)
	}

	// Slice to hold the results (pgxscan.Select will append to this)
	// Initialize as an empty slice
	attendances := []*models.Attendance{}

	err = pgxscan.Select(ctx, a.DB.Pool, &attendances, sql, args...) // Pass the address of the slice
	if err != nil {
		// pgxscan.Select returns nil error and an empty slice if no rows are found.
		// So, an error here typically indicates a problem with query execution or scanning errors during iteration.
		return nil, fmt.Errorf("GetAttendanceByCourse: failed to execute query or scan: %w", err) // Wrap the original error
	}

	// If no error occurred, attendances will contain the results (or be an empty slice if no rows matched)
	return attendances, nil
}

func (a *attendanceRepository) MarkAttendance(ctx context.Context, collegeID int, studentID, courseID int, lectureID int) (bool, error) {
	now := time.Now()
	// Truncate date for the 'date' column if you only store the date part
	attendanceDate := now.Truncate(24 * time.Hour)

	// This query attempts to insert a record.
	// If a record for the same student, course, lecture, and date already exists,
	// it updates the scanned_at timestamp. This is a common "upsert" pattern.
	query := a.DB.SQ.Insert(attendanceTable).
		Columns(
			"student_id",
			"course_id",
			"college_id",
			"lecture_id",
			"date",
			"status", // Initial status, e.g., "Present"
			"scanned_at",
		).
		Values(
			studentID,
			courseID,
			collegeID,
			lectureID,
			attendanceDate,
			"Present", // Default status when marked
			now,
		).
		Suffix(`ON CONFLICT (student_id, course_id, lecture_id, date, college_id)
              DO UPDATE SET scanned_at = EXCLUDED.scanned_at, status = EXCLUDED.status`) // Update scan time and status on conflict

	sql, args, err := query.ToSql()
	if err != nil {
		return false, fmt.Errorf("MarkAttendance: failed to build query: %w", err)
	}

	// Execute the query (Exec is used for INSERT/UPDATE/DELETE)
	commandTag, err := a.DB.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return false, fmt.Errorf("MarkAttendance: failed to execute query: %w", err)
	}

	// commandTag.RowsAffected() will be 1 if a row was inserted or updated.
	// It's a good check, but often just checking for nil error is sufficient for "success".
	// Given the bool return, let's return true if at least one row was affected.
	return commandTag.RowsAffected() > 0, nil
}

func (a *attendanceRepository) UpdateAttendance(ctx context.Context, collegeID int, studentID int, courseID int, lectureID int, status string) error {
	query := a.DB.SQ.Update(attendanceTable).From(attendanceTable).Set("status", status).Where(squirrel.Eq{
		"college_id": collegeID,
		"student_id": studentID,
		"course_id":  courseID,
		"lecture_id": lectureID,
	})
	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("failed to build query")
	}
	commandTag, err := a.DB.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("failed to query update attendance")
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("did not update attendance")
	}
	return nil
}

func (a *attendanceRepository) GetAttendanceStudentInCourse(ctx context.Context, collegeID int, studentID int, courseID int) ([]*models.Attendance, error) {
	attendances := []*models.Attendance{}
	query := a.DB.SQ.Select("id", // Use database column names matching struct tags
		"student_id",
		"course_id",
		"college_id",
		"date",
		"status",
		"scanned_at",
		"lecture_id").From(attendanceTable).Where(squirrel.Eq{
		"college_id": collegeID,
		"student_id": studentID,
		"course_id":  courseID,
	}).OrderBy("scanned_at ASC")
	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("unabel to build query ")
	}
	sqlErr := pgxscan.Select(ctx, a.DB.Pool, &attendances, sql, args...)
	if sqlErr != nil {
		return nil, fmt.Errorf("failed to execute getAttendanceStudentINCourse")
	}
	return attendances, nil
}

func (a *attendanceRepository) GetAttendanceStudent(ctx context.Context, collegeID int, studentID int) ([]*models.Attendance, error) {
	// Build the SELECT query for multiple rows
	query := a.DB.SQ.Select(
		"id", "student_id", "course_id", "college_id",
		"date", "status", "scanned_at", "lecture_id",
	).
		From(attendanceTable).
		Where(squirrel.Eq{
			"college_id": collegeID,
			"student_id": studentID,
		}).
		OrderBy("date ASC", "course_id ASC", "scanned_at ASC") // Optional: Add ordering

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("GetAttendanceStudent: failed to build query: %w", err)
	}

	attendances := []*models.Attendance{} // Slice to hold results

	// Use pgxscan.Select for multiple rows
	err = pgxscan.Select(ctx, a.DB.Pool, &attendances, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("GetAttendanceStudent: failed to execute query or scan: %w", err)
	}

	return attendances, nil // Returns slice (empty if no rows) and nil error on success
}

// GetAttendanceByLecture retrieves attendance records for a specific lecture.
func (a *attendanceRepository) GetAttendanceByLecture(ctx context.Context, collegeID int, lectureID int, courseID int) ([]*models.Attendance, error) {
	// Build the SELECT query for multiple rows
	query := a.DB.SQ.Select(
		"id", "student_id", "course_id", "college_id",
		"date", "status", "scanned_at", "lecture_id",
	).
		From(attendanceTable).
		Where(squirrel.Eq{
			"college_id": collegeID,
			"lecture_id": lectureID,
			"course_id":  courseID, // Include courseID as per the interface, even if lectureID might be globally unique
		}).
		OrderBy("student_id ASC", "scanned_at ASC") // Optional: Add ordering

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("GetAttendanceByLecture: failed to build query: %w", err)
	}

	attendances := []*models.Attendance{} // Slice to hold results

	// Use pgxscan.Select for multiple rows
	err = pgxscan.Select(ctx, a.DB.Pool, &attendances, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("GetAttendanceByLecture: failed to execute query or scan: %w", err)
	}

	return attendances, nil // Returns slice (empty if no rows) and nil error on success
}

// FreezeAttendance updates the status of all attendance records for a specific student to "Frozen".
// This is a simple example; actual freezing logic might be more complex (e.g., only for past dates).
func (a *attendanceRepository) FreezeAttendance(ctx context.Context, collegeID int, studentID int) error {
	// Build the UPDATE query
	query := a.DB.SQ.Update(attendanceTable).
		Set("status", "Frozen"). // Set the status to "Frozen"
		// Add other potential fields like freeze_date = now() if needed
		Where(squirrel.Eq{ // Identify the records to freeze
			"college_id": collegeID,
			"student_id": studentID,
		})
		// Optional: Add a condition to only freeze records that aren't already frozen or finalized
		// .Where(squirrel.NotEq{"status": "Frozen"})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("FreezeAttendance: failed to build query: %w", err)
	}

	// Execute the query
	commandTag, err := a.DB.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("FreezeAttendance: failed to execute query: %w", err)
	}

	// Optional: Check if any rows were affected. Freezing might affect 0 rows
	// if the student has no attendance records or they are already frozen.
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("unabel to freeze attendance")
	}
	return nil // Success
}
