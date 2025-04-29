package repository

import (
	"context"
	"fmt"  // Import fmt for better error wrapping
	"time" // Assuming models.Enrollment uses time.Time

	"eduhub/server/internal/models" // Your models package
	// Removed "errors" if only using fmt.Errorf
	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan"
	"github.com/jackc/pgx/v4" // For pgx.ErrNoRows
	// Assuming pgxpool is imported elsewhere in your repository package
	// "github.com/jackc/pgx/v4/pgxpool"
	// Note: pgxscan is not strictly needed for Create and Exists if
	// you are scanning only the ID or a single value (like '1' for Exists),
	// but you'd use it for Get/Find/Select methods on Enrollment if you add them.
)

type EnrollmentRepository interface {
	CreateEnrollment(ctx context.Context, enrollment *models.Enrollment) error
	IsStudentEnrolled(ctx context.Context, collegeID, studentID, courseID int) (bool, error)
	// Add other specific enrollment methods here (e.g., GetByID, FindByStudent, UpdateStatus)
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

func (e *enrollmentRepository) GetEnrollmentByID(ctx context.Context, id int) (*models.Enrollment, error) {
	// Build the SELECT query for a single row
	query := e.DB.SQ.Select(
		"id", "student_id", "course_id", "college_id",
		"enrollment_date", "status", "grade", "created_at", "updated_at",
	).
		From(enrollmentTable).
		Where(squirrel.Eq{"id": id}) // Filter by ID

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
			return nil, fmt.Errorf("GetEnrollmentByID: enrollment with ID %d not found", id) // Or a custom ErrNotFound
		}
		// Any other error during execution or scanning
		return nil, fmt.Errorf("GetEnrollmentByID: failed to execute query or scan: %w", err)
	}

	return enrollment, nil // Success
}

// FindEnrollmentsByStudent retrieves all enrollment records for a specific student in a college.
func (e *enrollmentRepository) FindEnrollmentsByStudent(ctx context.Context, collegeID int, studentID int) ([]*models.Enrollment, error) {
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
		OrderBy("enrollment_date DESC", "course_id ASC") // Optional: Add ordering

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

// UpdateEnrollmentStatus updates the status of a specific enrollment record by ID.
func (e *enrollmentRepository) UpdateEnrollmentStatus(ctx context.Context, id int, status string) error {
	// Update the UpdatedAt timestamp
	now := time.Now()

	// Build the UPDATE query
	query := e.DB.SQ.Update(enrollmentTable).
		Set("status", status).       // Set the new status
		Set("updated_at", now).      // Update the updated_at timestamp
		Where(squirrel.Eq{"id": id}) // Identify the record by ID

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
		return fmt.Errorf("UpdateEnrollmentStatus: no row updated for ID %d", id)
	}

	return nil // Success
}
