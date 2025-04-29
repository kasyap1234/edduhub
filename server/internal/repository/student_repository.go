package repository

import (
	"context"
	"fmt"
	"time" // Needed for time.Now() and time.Time fields

	"eduhub/server/internal/models" // Your models package

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan" // Using pgxscan for GET
	"github.com/jackc/pgx/v4"              // For pgx.ErrNoRows, CommandTag
	// Assuming DB struct uses this
)

// Ensure DB struct is defined elsewhere in this package
// type DB struct {
// 	Pool *pgxpool.Pool
// 	SQ   squirrel.StatementBuilderType
// }

// --- Updated models.Student struct (assuming these fields exist in your DB) ---

type StudentRepository interface {
	CreateStudent(ctx context.Context, student *models.Student) error
	GetStudentByRollNo(ctx context.Context, rollNo string) (*models.Student, error)
	GetStudentByID(ctx context.Context, collegeID int, studentID int) (*models.Student, error) // Note: studentID is the primary key 'id' here
	UpdateStudent(ctx context.Context, model *models.Student) error
	FreezeStudent(ctx context.Context, rollNo string) error   // Renamed param to match casing
	UnFreezeStudent(ctx context.Context, rollNo string) error // Renamed param to match casing
	FindByKratosID(ctx context.Context, kratosID string) (*models.Student, error)
	// Note: GetStudentByID signature is a bit unusual if ID is the primary key;
	// typically you only need the ID. Assuming you intend to filter by both ID and CollegeID.
}

// studentRepository now holds a direct reference to *DB
type studentRepository struct {
	DB *DB
}

// NewStudentRepository receives the *DB directly
func NewStudentRepository(db *DB) StudentRepository {
	return &studentRepository{
		DB: db,
	}
}

const studentTable = "students" // Define your table name

// CreateStudent inserts a new student record into the database.
func (s *studentRepository) CreateStudent(ctx context.Context, student *models.Student) error {
	// Set timestamps if they are zero-valued
	now := time.Now()
	if student.CreatedAt.IsZero() {
		student.CreatedAt = now
	}
	if student.UpdatedAt.IsZero() {
		student.UpdatedAt = now
	}

	// Build the INSERT query using squirrel
	query := s.DB.SQ.Insert(studentTable).
		Columns(
			// Include all relevant fields, including the new ones
			"user_id",
			"college_id",
			"kratos_identity_id",
			"enrollment_year",
			"roll_no",   // New field
			"is_active", // New field (set initial status, e.g., true)
			"created_at",
			"updated_at",
		).
		Values(
			student.UserID,
			student.CollegeID,
			student.KratosIdentityID,
			student.EnrollmentYear,
			student.RollNo,
			student.IsActive, // Use the field value
			student.CreatedAt,
			student.UpdatedAt,
		).
		Suffix("RETURNING id") // Assuming 'id' is auto-generated and you want it back

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("CreateStudent: failed to build query: %w", err)
	}

	// Execute the query and scan the returned ID back into the struct
	err = s.DB.Pool.QueryRow(ctx, sql, args...).Scan(&student.ID)
	if err != nil {
		return fmt.Errorf("CreateStudent: failed to execute query or scan ID: %w", err)
	}

	return nil // Success
}

// GetStudentByRollNo retrieves a student by their roll number.
func (s *studentRepository) GetStudentByRollNo(ctx context.Context, rollNo string) (*models.Student, error) {
	// Build the SELECT query for a single row
	query := s.DB.SQ.Select(
		"id", "user_id", "college_id", "kratos_identity_id",
		"enrollment_year", "roll_no", "is_active", "created_at", "updated_at",
	).
		From(studentTable).
		Where(squirrel.Eq{"roll_no": rollNo}) // Filter by roll_no

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("GetStudentByRollNo: failed to build query: %w", err)
	}

	student := &models.Student{} // Initialize the struct

	// Use pgxscan.Get for a single row result
	err = pgxscan.Get(ctx, s.DB.Pool, student, sql, args...)
	if err != nil {
		if err == pgx.ErrNoRows {
			// Return nil, nil for not found, consistent with your original FindByKratosID
			return nil, nil
		}
		// Any other error during execution or scanning
		return nil, fmt.Errorf("GetStudentByRollNo: failed to execute query or scan: %w", err)
	}

	return student, nil // Success
}

// GetStudentByID retrieves a student by their ID, filtered by college ID.
func (s *studentRepository) GetStudentByID(ctx context.Context, collegeID int, studentID int) (*models.Student, error) {
	// Build the SELECT query for a single row
	query := s.DB.SQ.Select(
		"id", "user_id", "college_id", "kratos_identity_id",
		"enrollment_year", "roll_no", "is_active", "created_at", "updated_at",
	).
		From(studentTable).
		Where(squirrel.Eq{"id": studentID, "college_id": collegeID}) // Filter by ID AND CollegeID

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("GetStudentByID: failed to build query: %w", err)
	}

	student := &models.Student{} // Initialize the struct

	// Use pgxscan.Get for a single row result
	err = pgxscan.Get(ctx, s.DB.Pool, student, sql, args...)
	if err != nil {
		if err == pgx.ErrNoRows {
			// Return nil, nil for not found
			return nil, nil
		}
		// Any other error during execution or scanning
		return nil, fmt.Errorf("GetStudentByID: failed to execute query or scan: %w", err)
	}

	return student, nil // Success
}

// UpdateStudent updates an existing student record.
func (s *studentRepository) UpdateStudent(ctx context.Context, model *models.Student) error {
	// Update the UpdatedAt timestamp
	model.UpdatedAt = time.Now()

	// Build the UPDATE query
	// Note: You typically don't update ID or CreatedAt this way
	query := s.DB.SQ.Update(studentTable).
		Set("user_id", model.UserID).
		Set("college_id", model.CollegeID).
		Set("kratos_identity_id", model.KratosIdentityID).
		Set("enrollment_year", model.EnrollmentYear).
		Set("roll_no", model.RollNo).       // Include new field
		Set("is_active", model.IsActive).   // Include new field
		Set("updated_at", model.UpdatedAt). // Update timestamp
		Where(squirrel.Eq{"id": model.ID})  // Identify the record by ID

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("UpdateStudent: failed to build query: %w", err)
	}

	// Execute the query
	commandTag, err := s.DB.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("UpdateStudent: failed to execute query: %w", err)
	}

	// Optional: Check if a row was actually updated
	if commandTag.RowsAffected() == 0 {
		// You might want to return a specific error here if the ID wasn't found
		return fmt.Errorf("UpdateStudent: no row updated for ID %d", model.ID)
	}

	return nil // Success
}

// FreezeStudent sets the IsActive status of a student to false based on their roll number.
func (s *studentRepository) FreezeStudent(ctx context.Context, rollNo string) error {
	// Build the UPDATE query
	now := time.Now()
	query := s.DB.SQ.Update(studentTable).
		Set("is_active", false).              // Set status to false
		Set("updated_at", now).               // Update timestamp
		Where(squirrel.Eq{"roll_no": rollNo}) // Identify the student by roll_no

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("FreezeStudent: failed to build query: %w", err)
	}

	// Execute the query
	commandTag, err := s.DB.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("FreezeStudent: failed to execute query: %w", err)
	}

	// Optional: Check if a row was actually updated (i.e., roll number found)
	if commandTag.RowsAffected() == 0 {
		// You might want to return an error if the student wasn't found by roll number
		// return fmt.Errorf("FreezeStudent: student with roll number %s not found", rollNo)
		// Or, if freezing an already frozen student is okay, just proceed.
	}

	return nil // Success
}

// UnFreezeStudent sets the IsActive status of a student to true based on their roll number.
func (s *studentRepository) UnFreezeStudent(ctx context.Context, rollNo string) error {
	// Build the UPDATE query
	now := time.Now()
	query := s.DB.SQ.Update(studentTable).
		Set("is_active", true).               // Set status to true
		Set("updated_at", now).               // Update timestamp
		Where(squirrel.Eq{"roll_no": rollNo}) // Identify the student by roll_no

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("UnFreezeStudent: failed to build query: %w", err)
	}

	// Execute the query
	commandTag, err := s.DB.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("UnFreezeStudent: failed to execute query: %w", err)
	}

	// Optional: Check if a row was actually updated
	if commandTag.RowsAffected() == 0 {
		// return fmt.Errorf("UnFreezeStudent: student with roll number %s not found", rollNo)
	}

	return nil // Success
}

// FindByKratosID retrieves a student record by their Kratos identity ID.
func (s *studentRepository) FindByKratosID(ctx context.Context, kratosID string) (*models.Student, error) {
	// Build the SELECT query for a single row
	query := s.DB.SQ.Select(
		"id", "user_id", "college_id", "kratos_identity_id",
		"enrollment_year", "roll_no", "is_active", "created_at", "updated_at",
	).
		From(studentTable).
		Where(squirrel.Eq{"kratos_identity_id": kratosID}) // Filter by kratos_identity_id

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("FindByKratosID: failed to build query: %w", err)
	}

	student := &models.Student{} // Initialize the struct

	// Use pgxscan.Get for a single row result
	err = pgxscan.Get(ctx, s.DB.Pool, student, sql, args...)
	if err != nil {
		if err == pgx.ErrNoRows {
			// Return nil, nil for not found, consistent with your original code
			return nil, nil
		}
		// Any other error during execution or scanning
		return nil, fmt.Errorf("FindByKratosID: failed to execute query or scan: %w", err)
	}

	return student, nil // Success
}
