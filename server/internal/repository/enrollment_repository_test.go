package repository

import (
	"context"
	"errors"
	"regexp"
	"testing"
	"time"

	"eduhub/server/internal/models"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupEnrollmentTest(t *testing.T) (pgxmock.PgxPoolIface, *DB, EnrollmentRepository, context.Context) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)

	db := &DB{
		Pool: mock,
		SQ:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}

	// Assuming NewEnrollmentRepository takes *DB
	repo := NewEnrollmentRepository(db)
	ctx := context.Background()

	return mock, db, repo, ctx
}

func TestCreateEnrollment(t *testing.T) {
	mock, _, repo, ctx := setupEnrollmentTest(t)
	defer mock.Close()

	enrollment := &models.Enrollment{
		StudentID:      101,
		CourseID:       202,
		CollegeID:      1,
		EnrollmentDate: time.Now().Truncate(24 * time.Hour),
		Status:         "Enrolled",
		Grade:          "N/A",
		// CreatedAt and UpdatedAt will be set by the method if zero
	}
	expectedID := 50

	sqlRegex := `INSERT INTO enrollments \(student_id,course_id,college_id,enrollment_date,status,grade,created_at,updated_at\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7,\$8\) RETURNING id`

	mock.ExpectQuery(regexp.QuoteMeta(sqlRegex)).
		WithArgs(
			enrollment.StudentID,
			enrollment.CourseID,
			enrollment.CollegeID,
			enrollment.EnrollmentDate,
			enrollment.Status,
			enrollment.Grade,
			pgxmock.AnyArg(), // CreatedAt
			pgxmock.AnyArg(), // UpdatedAt
		).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(expectedID))

	err := repo.CreateEnrollment(ctx, enrollment)

	assert.NoError(t, err)
	assert.Equal(t, expectedID, enrollment.ID)
	assert.False(t, enrollment.CreatedAt.IsZero())
	assert.False(t, enrollment.UpdatedAt.IsZero())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateEnrollment_Error(t *testing.T) {
	mock, _, repo, ctx := setupEnrollmentTest(t)
	defer mock.Close()

	enrollment := &models.Enrollment{
		StudentID: 101,
		CourseID:  202,
		CollegeID: 1,
		Status:    "Enrolled",
	}
	dbError := errors.New("db connection failed")
	sqlRegex := `INSERT INTO enrollments \(student_id,course_id,college_id,enrollment_date,status,grade,created_at,updated_at\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7,\$8\) RETURNING id`

	mock.ExpectQuery(regexp.QuoteMeta(sqlRegex)).
		WithArgs(enrollment.StudentID, enrollment.CourseID, enrollment.CollegeID, pgxmock.AnyArg(), enrollment.Status, pgxmock.AnyArg(), pgxmock.AnyArg(), pgxmock.AnyArg()).
		WillReturnError(dbError)

	err := repo.CreateEnrollment(ctx, enrollment)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to execute query or scan ID")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestIsStudentEnrolled_True(t *testing.T) {
	mock, _, repo, ctx := setupEnrollmentTest(t)
	defer mock.Close()

	collegeID := 1
	studentID := 101
	courseID := 202
	sqlRegex := `SELECT 1 FROM enrollments WHERE college_id = \$1 AND course_id = \$2 AND student_id = \$3`

	mock.ExpectQuery(regexp.QuoteMeta(sqlRegex)).
		WithArgs(collegeID, courseID, studentID). // Order matters based on squirrel's output
		WillReturnRows(pgxmock.NewRows([]string{"1"}).AddRow(1))

	exists, err := repo.IsStudentEnrolled(ctx, collegeID, studentID, courseID)

	assert.NoError(t, err)
	assert.True(t, exists)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestIsStudentEnrolled_False(t *testing.T) {
	mock, _, repo, ctx := setupEnrollmentTest(t)
	defer mock.Close()

	collegeID := 1
	studentID := 101
	courseID := 202
	sqlRegex := `SELECT 1 FROM enrollments WHERE college_id = \$1 AND course_id = \$2 AND student_id = \$3`

	mock.ExpectQuery(regexp.QuoteMeta(sqlRegex)).
		WithArgs(collegeID, courseID, studentID).
		WillReturnError(pgx.ErrNoRows)

	exists, err := repo.IsStudentEnrolled(ctx, collegeID, studentID, courseID)

	assert.NoError(t, err)
	assert.False(t, exists)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestIsStudentEnrolled_Error(t *testing.T) {
	mock, _, repo, ctx := setupEnrollmentTest(t)
	defer mock.Close()

	collegeID := 1
	studentID := 101
	courseID := 202
	dbError := errors.New("query failed")
	sqlRegex := `SELECT 1 FROM enrollments WHERE college_id = \$1 AND course_id = \$2 AND student_id = \$3`

	mock.ExpectQuery(regexp.QuoteMeta(sqlRegex)).
		WithArgs(collegeID, courseID, studentID).
		WillReturnError(dbError)

	exists, err := repo.IsStudentEnrolled(ctx, collegeID, studentID, courseID)

	assert.Error(t, err)
	assert.False(t, exists)
	assert.Contains(t, err.Error(), "failed to execute query")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// --- Tests for GetEnrollmentByID ---

func TestGetEnrollmentByID(t *testing.T) {
	mock, _, repo, ctx := setupEnrollmentTest(t)
	defer mock.Close()

	enrollmentID := 50
	expectedEnrollment := &models.Enrollment{
		ID:             enrollmentID,
		StudentID:      101,
		CourseID:       202,
		CollegeID:      1,
		EnrollmentDate: time.Now().Add(-10 * 24 * time.Hour).Truncate(24 * time.Hour),
		Status:         "Completed",
		Grade:          "A",
		CreatedAt:      time.Now().Add(-11 * 24 * time.Hour),
		UpdatedAt:      time.Now().Add(-1 * 24 * time.Hour),
	}

	sqlRegex := `SELECT id, student_id, course_id, college_id, enrollment_date, status, grade, created_at, updated_at FROM enrollments WHERE id = \$1`
	rows := pgxmock.NewRows([]string{"id", "student_id", "course_id", "college_id", "enrollment_date", "status", "grade", "created_at", "updated_at"}).
		AddRow(expectedEnrollment.ID, expectedEnrollment.StudentID, expectedEnrollment.CourseID, expectedEnrollment.CollegeID, expectedEnrollment.EnrollmentDate, expectedEnrollment.Status, expectedEnrollment.Grade, expectedEnrollment.CreatedAt, expectedEnrollment.UpdatedAt)

	// pgxscan.Get uses QueryRow which expects a single row result.
	mock.ExpectQuery(regexp.QuoteMeta(sqlRegex)).
		WithArgs(enrollmentID).
		WillReturnRows(rows)

	enrollment, err := repo.GetEnrollmentByID(ctx, enrollmentID)

	assert.NoError(t, err)
	// Compare fields individually or use assert.EqualValues for structs if appropriate
	assert.Equal(t, expectedEnrollment.ID, enrollment.ID)
	assert.Equal(t, expectedEnrollment.StudentID, enrollment.StudentID)
	assert.Equal(t, expectedEnrollment.Status, enrollment.Status)
	// ... compare other fields ...
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestGetEnrollmentByID_NotFound(t *testing.T) {
	mock, _, repo, ctx := setupEnrollmentTest(t)
	defer mock.Close()

	enrollmentID := 999
	sqlRegex := `SELECT id, student_id, course_id, college_id, enrollment_date, status, grade, created_at, updated_at FROM enrollments WHERE id = \$1`

	mock.ExpectQuery(regexp.QuoteMeta(sqlRegex)).
		WithArgs(enrollmentID).
		WillReturnError(pgx.ErrNoRows) // Simulate not found

	enrollment, err := repo.GetEnrollmentByID(ctx, enrollmentID)

	assert.Error(t, err)
	assert.Nil(t, enrollment)
	assert.Contains(t, err.Error(), "not found") // Check the specific error message from the repo
	assert.NoError(t, mock.ExpectationsWereMet())
}

// --- Tests for FindEnrollmentsByStudent (Example) ---
// Add tests for FindEnrollmentsByStudent similar to GetAttendanceByCourse

// --- Tests for UpdateEnrollmentStatus ---

func TestUpdateEnrollmentStatus(t *testing.T) {
	mock, _, repo, ctx := setupEnrollmentTest(t)
	defer mock.Close()

	enrollmentID := 50
	newStatus := "Withdrawn"

	sqlRegex := `UPDATE enrollments SET status = \$1, updated_at = \$2 WHERE id = \$3`

	mock.ExpectExec(regexp.QuoteMeta(sqlRegex)).
		WithArgs(newStatus, pgxmock.AnyArg(), enrollmentID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1)) // Simulate 1 row affected

	err := repo.UpdateEnrollmentStatus(ctx, enrollmentID, newStatus)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateEnrollmentStatus_NotFound(t *testing.T) {
	mock, _, repo, ctx := setupEnrollmentTest(t)
	defer mock.Close()

	enrollmentID := 999
	newStatus := "Withdrawn"
	sqlRegex := `UPDATE enrollments SET status = \$1, updated_at = \$2 WHERE id = \$3`

	mock.ExpectExec(regexp.QuoteMeta(sqlRegex)).
		WithArgs(newStatus, pgxmock.AnyArg(), enrollmentID).
		WillReturnResult(pgxmock.NewResult("UPDATE", 0)) // Simulate 0 rows affected

	err := repo.UpdateEnrollmentStatus(ctx, enrollmentID, newStatus)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no row updated") // Check the specific error message
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateEnrollmentStatus_Error(t *testing.T) {
	mock, _, repo, ctx := setupEnrollmentTest(t)
	defer mock.Close()

	enrollmentID := 50
	newStatus := "Withdrawn"
	dbError := errors.New("update failed")
	sqlRegex := `UPDATE enrollments SET status = \$1, updated_at = \$2 WHERE id = \$3`

	mock.ExpectExec(regexp.QuoteMeta(sqlRegex)).
		WithArgs(newStatus, pgxmock.AnyArg(), enrollmentID).
		WillReturnError(dbError)

	err := repo.UpdateEnrollmentStatus(ctx, enrollmentID, newStatus)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to execute query")
	assert.NoError(t, mock.ExpectationsWereMet())
}