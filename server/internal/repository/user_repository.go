package repository

import (
	"context"
	"fmt"  // Import fmt for better error wrapping
	"time" // Needed for time.Now() and time.Time fields

	"eduhub/server/internal/models" // Your models package

	"github.com/Masterminds/squirrel"
	// Using pgxscan for GET
	// For pgx.ErrNoRows, CommandTag
	// Assuming DB struct uses this
)

// Ensure DB struct is defined elsewhere in this package
// type DB struct {
// 	Pool *pgxpool.Pool
// 	SQ   squirrel.StatementBuilderType
// }

// --- Updated models.User struct (assuming RollNo field exists in your DB) ---
// type User struct {
//  ID               int       `db:"id" json:"id"`
//  Name             string    `db:"name" json:"name"`
//  Role             string    `db:"role" json:"role"`
//  Email            string    `db:"email" json:"email"`
//  KratosIdentityID string    `db:"kratos_identity_id" json:"kratos_identity_id"`
//  IsActive         bool      `db:"is_active" json:"is_active"`
//  RollNo           string    `db:"roll_no" json:"roll_no"` // Added this field
//  CreatedAt        time.Time `db:"created_at" json:"created_at"`
//  UpdatedAt        time.Time `db:"updated_at" json:"updated_at"`

//  // Relations - not stored in DB (add db:"-" tag)
//  // Student *Student `db:"-" json:"student,omitempty"`
// }

type UserRepository interface {
	CreateUser(ctx context.Context, user *models.User) error

	UpdateUser(ctx context.Context, user *models.User) error
	FreezeUser(ctx context.Context, rollNo string) error // Renamed param to match casing
	DeleteUser(ctx context.Context, rollNo string) error // Renamed param to match casing
	// Note: methods like GetUserByID would be useful here too
}

// userRepository now holds a direct reference to *DB
type userRepository struct {
	DB *DB
}

// NewUserRepository receives the *DB directly
func NewUserRepository(db *DB) UserRepository {
	return &userRepository{
		DB: db,
	}
}

const userTable = "users" // Define your table name

// CreateUser inserts a new user record into the database.
func (u *userRepository) CreateUser(ctx context.Context, user *models.User) error {
	// Set timestamps if they are zero-valued
	now := time.Now()
	if user.CreatedAt.IsZero() {
		user.CreatedAt = now
	}
	if user.UpdatedAt.IsZero() {
		user.UpdatedAt = now
	}

	// Build the INSERT query using squirrel
	query := u.DB.SQ.Insert(userTable).
		Columns(
			"name",
			"role",
			"email",
			"kratos_identity_id",
			"is_active",
			"roll_no", // Include the new field
			"created_at",
			"updated_at",
		).
		Values(
			user.Name,
			user.Role,
			user.Email,
			user.KratosIdentityID,
			user.IsActive,

			user.CreatedAt,
			user.UpdatedAt,
		).
		Suffix("RETURNING id") // Assuming 'id' is auto-generated and you want it back

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("CreateUser: failed to build query: %w", err)
	}

	// Execute the query and scan the returned ID back into the struct
	err = u.DB.Pool.QueryRow(ctx, sql, args...).Scan(&user.ID)
	if err != nil {
		return fmt.Errorf("CreateUser: failed to execute query or scan ID: %w", err)
	}

	return nil // Success
}

// UpdateUser updates an existing user record.
func (u *userRepository) UpdateUser(ctx context.Context, model *models.User) error {
	// Update the UpdatedAt timestamp
	model.UpdatedAt = time.Now()

	// Build the UPDATE query
	// Note: You typically don't update ID or CreatedAt this way
	query := u.DB.SQ.Update(userTable).
		Set("name", model.Name).
		Set("role", model.Role).
		Set("email", model.Email).
		Set("kratos_identity_id", model.KratosIdentityID).
		Set("is_active", model.IsActive).
		// Include the roll_no field
		Set("updated_at", model.UpdatedAt). // Update timestamp
		Where(squirrel.Eq{"id": model.ID})  // Identify the record by ID

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("UpdateUser: failed to build query: %w", err)
	}

	// Execute the query
	commandTag, err := u.DB.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("UpdateUser: failed to execute query: %w", err)
	}

	// Optional: Check if a row was actually updated
	if commandTag.RowsAffected() == 0 {
		// You might want to return a specific error here if the ID wasn't found
		return fmt.Errorf("UpdateUser: no row updated for ID %d", model.ID)
	}

	return nil // Success
}

// FreezeUser sets the IsActive status of a user to false based on their roll number.
// This implementation updates directly by roll_no without fetching first.
func (u *userRepository) FreezeUser(ctx context.Context, rollNo string) error {
	// Build the UPDATE query
	now := time.Now()
	query := u.DB.SQ.Update(userTable).
		Set("is_active", false).              // Set status to false
		Set("updated_at", now).               // Update timestamp
		Where(squirrel.Eq{"roll_no": rollNo}) // Identify the user by roll_no

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("FreezeUser: failed to build query: %w", err)
	}

	// Execute the query
	commandTag, err := u.DB.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("FreezeUser: failed to execute query: %w", err)
	}

	// Optional: Check if a row was actually affected (i.e., roll number found).
	// If freezing an already frozen user is okay, this check might not be strictly necessary
	// depending on whether you need to know if a change *actually* happened.
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("FreezeUser: user with roll number %s not found or already frozen", rollNo)
	}

	return nil // Success
}

// DeleteUser deletes a user record based on their roll number.
// This implementation deletes directly by roll_no without fetching first.
func (u *userRepository) DeleteUser(ctx context.Context, rollNo string) error {
	// Build the DELETE query
	query := u.DB.SQ.Delete(userTable).
		Where(squirrel.Eq{"roll_no": rollNo}) // Identify the user by roll_no

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("DeleteUser: failed to build query: %w", err)
	}

	// Execute the query
	commandTag, err := u.DB.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("DeleteUser: failed to execute query: %w", err)
	}

	// Optional: Check if a row was actually deleted
	if commandTag.RowsAffected() == 0 {
		// You might want to return an error if the user wasn't found by roll number
		return fmt.Errorf("DeleteUser: user with roll number %s not found", rollNo)
	}

	return nil // Success
}
