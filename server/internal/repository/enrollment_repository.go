package repository

import (
	"context"
	"fmt"
	"time"

	"eduhub/server/internal/models"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4"
)

type EnrollmentRepository interface {
	CreateEnrollment(ctx context.Context, enrollment *models.Enrollment) error
	IsStudentEnrolled(ctx context.Context, collegeID int, studentID int, courseID int) (bool, error)
	GetEnrollmentByID(ctx context.Context, collegeID int, enrollmentID int) (*models.Enrollment, error) // Added collegeID for scoping
	UpdateEnrollment(ctx context.Context, enrollment *models.Enrollment) error
	UpdateEnrollmentStatus(ctx context.Context, collegeID int, enrollmentID int, status string) error // Added collegeID for scoping
	DeleteEnrollment(ctx context.Context, collegeID int, enrollmentID int) error

	// Find methods with pagination
	FindEnrollmentsByStudent(ctx context.Context, collegeID int, studentID int, limit, offset uint64) ([]*models.Enrollment, error)
	FindEnrollmentsByCourse(ctx context.Context, collegeID int, courseID int, limit, offset uint64) ([]*models.Enrollment, error)
	FindEnrollmentsByCollege(ctx context.Context, collegeID int, limit, offset uint64) ([]*models.Enrollment, error)

	// Count methods
	CountEnrollmentsByStudent(ctx context.Context, collegeID int, studentID int) (int, error)
	CountEnrollmentsByCourse(ctx context.Context, collegeID int, courseID int) (int, error)
	CountEnrollmentsByCollege(ctx context.Context, collegeID int) (int, error)
}

// enrollmentRepository now holds a direct reference to *DB
type enrollmentRepository struct {
	DB *DB
}

// NewEnrollmentRepository receives the *DB directly
func NewEnrollmentRepository(db *DB) EnrollmentRepository {
	return &enrollmentRepository{
		DB: db,
	}
}

const enrollmentTable = "enrollments" // Define your table name

// CreateEnrollment inserts a new enrollment record into the database.
func (e *enrollmentRepository) CreateEnrollment(ctx context.Context, enrollment *models.Enrollment) error {
	// Set timestamps if they are zero-valued
	now := time.Now()
	if enrollment.CreatedAt.IsZero() {
		enrollment.CreatedAt = now
	}
	if enrollment.UpdatedAt.IsZero() {
		enrollment.UpdatedAt = now
	}

	// Build the INSERT query using squirrel
	query := e.DB.SQ.Insert(enrollmentTable).
		Columns(
			"student_id",
			"course_id",
			"college_id", // Include college_id based on your struct and queries
			"enrollment_date",
			"status",
			"grade",
			"created_at",
			"updated_at",
		).
		Values(
			enrollment.StudentID,
			enrollment.CourseID,
			enrollment.CollegeID, // Use the field from the struct
			enrollment.EnrollmentDate,
			enrollment.Status,
			enrollment.Grade,
			enrollment.CreatedAt,
			enrollment.UpdatedAt,
		).
		Suffix("RETURNING id") // Assuming 'id' is auto-generated and you want it back

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("CreateEnrollment: failed to build query: %w", err)
	}

	// Execute the query and scan the returned ID back into the struct
	err = e.DB.Pool.QueryRow(ctx, sql, args...).Scan(&enrollment.ID)
	if err != nil {
		return fmt.Errorf("CreateEnrollment: failed to execute query or scan ID: %w", err)
	}

	return nil // Success
}

// IsStudentEnrolled checks if a student is enrolled in a specific course within a college.
func (e *enrollmentRepository) IsStudentEnrolled(ctx context.Context, collegeID int, studentID int, courseID int) (bool, error) {
	// Build a SELECT 1 query - efficient for checking existence
	query := e.DB.SQ.Select("1").
		From(enrollmentTable).
		Where(squirrel.Eq{ // Use squirrel.Eq for type-safe filtering
			"college_id": collegeID,
			"student_id": studentID,
			"course_id":  courseID,
		})

	sql, args, err := query.ToSql()
	if err != nil {
		return false, fmt.Errorf("IsStudentEnrolled: failed to build query: %w", err)
	}

	var exists int // Dummy variable to scan the '1' into

	// Execute the query expecting at most one row
	err = e.DB.Pool.QueryRow(ctx, sql, args...).Scan(&exists)
	if err != nil {
		if err == pgx.ErrNoRows {
			// If ErrNoRows is returned, the record does NOT exist
			return false, nil
		}
		// Any other error indicates a problem with query execution
		return false, fmt.Errorf("IsStudentEnrolled: failed to execute query: %w", err)
	}

	// If Scan succeeded without ErrNoRows, it means a row was returned ('1'), so the record exists
	return true, nil // Record exists
}

// GetEnrollmentByID retrieves a specific enrollment by its ID, scoped by collegeID.
func (e *enrollmentRepository) GetEnrollmentByID(ctx context.Context, collegeID int, enrollmentID int) (*models.Enrollment, error) {
	// Build the SELECT query for a single row
	query := e.DB.SQ.Select(
		"id", "student_id", "course_id", "college_id",
		"enrollment_date", "status", "grade", "created_at", "updated_at",
	).
		From(enrollmentTable).
		Where(squirrel.Eq{ // Filter by ID and CollegeID
			"id":         enrollmentID,
			"college_id": collegeID,
		})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("GetEnrollmentByID: failed to build query: %w", err)
	}

	enrollment := &models.Enrollment{} // Initialize the struct

	// Use pgxscan.Get for a single row result
	err = pgxscan.Get(ctx, e.DB.Pool, enrollment, sql, args...)
	if err != nil {
		if err == pgx.ErrNoRows {
			// Return nil and a specific error if no record is found
			return nil, fmt.Errorf("GetEnrollmentByID: enrollment with ID %d not found for college ID %d", enrollmentID, collegeID) // Or a custom ErrNotFound
		}
		// Any other error during execution or scanning
		return nil, fmt.Errorf("GetEnrollmentByID: failed to execute query or scan: %w", err)
	}

	return enrollment, nil // Success
}

// FindEnrollmentsByStudent retrieves all enrollment records for a specific student in a college.
func (e *enrollmentRepository) FindEnrollmentsByStudent(ctx context.Context, collegeID int, studentID int, limit, offset uint64) ([]*models.Enrollment, error) {
	// Build the SELECT query for multiple rows
	query := e.DB.SQ.Select(
		"id", "student_id", "course_id", "college_id",
		"enrollment_date", "status", "grade", "created_at", "updated_at",
	).
		From(enrollmentTable).
		Where(squirrel.Eq{
			"college_id": collegeID,
			"student_id": studentID,
		}).
		OrderBy("enrollment_date DESC", "course_id ASC"). // Optional: Add ordering
		Limit(limit).                                     // Add LIMIT for pagination
		Offset(offset)                                    // Add OFFSET for pagination

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("FindEnrollmentsByStudent: failed to build query: %w", err)
	}

	enrollments := []*models.Enrollment{} // Slice to hold results

	// Use pgxscan.Select for multiple rows
	err = pgxscan.Select(ctx, e.DB.Pool, &enrollments, sql, args...) // Pass the address of the slice
	if err != nil {
		// Select returns nil error and an empty slice if no rows are found.
		// Check for errors indicating failure to execute or scan.
		return nil, fmt.Errorf("FindEnrollmentsByStudent: failed to execute query or scan: %w", err)
	}

	return enrollments, nil // Returns slice (empty if no rows) and nil error on success
}

// FindEnrollmentsByCourse retrieves all enrollment records for a specific course in a college with pagination.
func (e *enrollmentRepository) FindEnrollmentsByCourse(ctx context.Context, collegeID int, courseID int, limit, offset uint64) ([]*models.Enrollment, error) {
	query := e.DB.SQ.Select(
		"id", "student_id", "course_id", "college_id",
		"enrollment_date", "status", "grade", "created_at", "updated_at",
	).
		From(enrollmentTable).
		Where(squirrel.Eq{
			"college_id": collegeID,
			"course_id":  courseID,
		}).
		OrderBy("student_id ASC", "enrollment_date DESC"). // Order by student then date
		Limit(limit).
		Offset(offset)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("FindEnrollmentsByCourse: failed to build query: %w", err)
	}

	enrollments := []*models.Enrollment{}
	err = pgxscan.Select(ctx, e.DB.Pool, &enrollments, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("FindEnrollmentsByCourse: failed to execute query or scan: %w", err)
	}

	return enrollments, nil
}

// FindEnrollmentsByCollege retrieves all enrollment records for a specific college with pagination.
func (e *enrollmentRepository) FindEnrollmentsByCollege(ctx context.Context, collegeID int, limit, offset uint64) ([]*models.Enrollment, error) {
	query := e.DB.SQ.Select(
		"id", "student_id", "course_id", "college_id",
		"enrollment_date", "status", "grade", "created_at", "updated_at",
	).
		From(enrollmentTable).
		Where(squirrel.Eq{
			"college_id": collegeID,
		}).
		OrderBy("course_id ASC", "student_id ASC", "enrollment_date DESC"). // Order by course, student, then date
		Limit(limit).
		Offset(offset)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("FindEnrollmentsByCollege: failed to build query: %w", err)
	}

	enrollments := []*models.Enrollment{}
	err = pgxscan.Select(ctx, e.DB.Pool, &enrollments, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("FindEnrollmentsByCollege: failed to execute query or scan: %w", err)
	}

	return enrollments, nil
}

// CountEnrollmentsByStudent counts the total number of enrollments for a specific student in a college.
func (e *enrollmentRepository) CountEnrollmentsByStudent(ctx context.Context, collegeID int, studentID int) (int, error) {
	return e.countEnrollments(ctx, squirrel.Eq{"college_id": collegeID, "student_id": studentID})
}

// CountEnrollmentsByCourse counts the total number of enrollments for a specific course in a college.
func (e *enrollmentRepository) CountEnrollmentsByCourse(ctx context.Context, collegeID int, courseID int) (int, error) {
	return e.countEnrollments(ctx, squirrel.Eq{"college_id": collegeID, "course_id": courseID})
}

// CountEnrollmentsByCollege counts the total number of enrollments within a specific college.
func (e *enrollmentRepository) CountEnrollmentsByCollege(ctx context.Context, collegeID int) (int, error) {
	return e.countEnrollments(ctx, squirrel.Eq{"college_id": collegeID})
}

// countEnrollments is a helper function for counting based on conditions.
func (e *enrollmentRepository) countEnrollments(ctx context.Context, whereClause squirrel.Sqlizer) (int, error) {
	query := e.DB.SQ.Select("COUNT(*)").From(enrollmentTable).Where(whereClause)

	sql, args, err := query.ToSql()
	if err != nil {
		return 0, fmt.Errorf("countEnrollments: failed to build query: %w", err)
	}

	var count int
	err = e.DB.Pool.QueryRow(ctx, sql, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("countEnrollments: failed to execute query or scan: %w", err)
	}

	return count, nil
}

// UpdateEnrollmentStatus updates the status of a specific enrollment record by ID.
func (e *enrollmentRepository) UpdateEnrollmentStatus(ctx context.Context, collegeID int, enrollmentID int, status string) error {
	// Update the UpdatedAt timestamp
	now := time.Now()

	// Build the UPDATE query
	query := e.DB.SQ.Update(enrollmentTable).
		Set("status", status).       // Set the new status
		Set("updated_at", now).      // Update the updated_at timestamp
		Where(squirrel.Eq{ // Identify the record by ID and CollegeID
			"id":         enrollmentID,
			"college_id": collegeID,
		})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("UpdateEnrollmentStatus: failed to build query: %w", err)
	}

	// Execute the query (Exec is used for INSERT/UPDATE/DELETE)
	commandTag, err := e.DB.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("UpdateEnrollmentStatus: failed to execute query: %w", err)
	}

	// Optional: Check if a row was actually updated
	if commandTag.RowsAffected() == 0 {
		// You might want to return a specific error here if the ID wasn't found
		return fmt.Errorf("UpdateEnrollmentStatus: no enrollment found with ID %d for college ID %d, or status unchanged", enrollmentID, collegeID)
	}

	return nil // Success
}

// DeleteEnrollment removes an enrollment record by its ID, scoped by collegeID.
func (e *enrollmentRepository) DeleteEnrollment(ctx context.Context, collegeID, enrollmentID int) error {
	query := e.DB.SQ.Delete(enrollmentTable).Where(squirrel.Eq{
		"id":         enrollmentID,
		"college_id": collegeID, // Ensure deletion is scoped
	})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("DeleteEnrollment: failed to build query: %w", err)
	}

	commandTag, err := e.DB.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("DeleteEnrollment: failed to execute query: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("DeleteEnrollment: no enrollment found with ID %d for college ID %d, or already deleted", enrollmentID, collegeID)
	}

	return nil
}

// UpdateEnrollment updates mutable fields of an existing enrollment record.
func (e *enrollmentRepository) UpdateEnrollment(ctx context.Context, enrollment *models.Enrollment) error {
	enrollment.UpdatedAt = time.Now()
	query := e.DB.SQ.Update(enrollmentTable).
		Set("enrollment_date", enrollment.EnrollmentDate).
		Set("status", enrollment.Status).
		Set("grade", enrollment.Grade).
		Set("updated_at", enrollment.UpdatedAt). // Corrected typo: updated_at
		Where(squirrel.Eq{"id": enrollment.ID, "college_id": enrollment.CollegeID}) // Ensure update is scoped

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("UpdateEnrollment: failed to build query: %w", err)
	}

	commandTag, err := e.DB.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("UpdateEnrollment: failed to execute query: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("UpdateEnrollment: no enrollment found with ID %d for college ID %d, or no changes made", enrollment.ID, enrollment.CollegeID)
	}
	return nil
}
