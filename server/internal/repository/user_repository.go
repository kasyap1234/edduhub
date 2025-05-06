package repository

import (
	"context"
	"errors"
	"fmt"  // Import fmt for better error wrapping
	"time" // Needed for time.Now() and time.Time fields

	"eduhub/server/internal/models" // Your models package

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan" // Using pgxscan for GET
	"github.com/jackc/pgx/v4"              // For pgx.ErrNoRows, CommandTag
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
	FreezeUserByID(ctx context.Context, userID int) error // Changed to operate on ID
	DeleteUserByID(ctx context.Context, userID int) error // Changed to operate on ID
	GetUserByID(ctx context.Context, userID int) (*models.User, error)
	GetUserByKratosID(ctx context.Context, kratosID string) (*models.User, error)
	UnFreezeUserByID(ctx context.Context, userID int) error

	FindAllUsers(ctx context.Context, limit, offset uint64) ([]*models.User, error)
	CountUsers(ctx context.Context) (int, error)
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
			"created_at", // Removed roll_no
			"updated_at",
		).
		Values(
			user.Name,
			user.Role,
			user.Email,
			user.KratosIdentityID,
			user.IsActive,
			// user.RollNo, // Removed RollNo value
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

// GetUserByID retrieves a user by their primary ID.
func (u *userRepository) GetUserByID(ctx context.Context, userID int) (*models.User, error) {
	query := u.DB.SQ.Select(
		"id", "name", "role", "email", "kratos_identity_id", "is_active", "created_at", "updated_at",
	).
		From(userTable).
		Where(squirrel.Eq{"id": userID})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("GetUserByID: failed to build query: %w", err)
	}

	user := &models.User{}
	err = pgxscan.Get(ctx, u.DB.Pool, user, sql, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, fmt.Errorf("GetUserByID: user with ID %d not found", userID)
		}
		return nil, fmt.Errorf("GetUserByID: failed to execute query or scan: %w", err)
	}

	return user, nil
}

// GetUserByKratosID retrieves a user by their Kratos identity ID.
func (u *userRepository) GetUserByKratosID(ctx context.Context, kratosID string) (*models.User, error) {
	query := u.DB.SQ.Select(
		"id", "name", "role", "email", "kratos_identity_id", "is_active", "created_at", "updated_at",
	).
		From(userTable).
		Where(squirrel.Eq{"kratos_identity_id": kratosID})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("GetUserByKratosID: failed to build query: %w", err)
	}

	user := &models.User{}
	err = pgxscan.Get(ctx, u.DB.Pool, user, sql, args...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			// Return nil, nil consistent with FindByKratosID in student repo if preferred
			return nil, fmt.Errorf("GetUserByKratosID: user with Kratos ID %s not found", kratosID)
		}
		return nil, fmt.Errorf("GetUserByKratosID: failed to execute query or scan: %w", err)
	}

	return user, nil
}

// UnFreezeUserByID sets the IsActive status of a user to true based on their ID.
func (u *userRepository) UnFreezeUserByID(ctx context.Context, userID int) error {
	now := time.Now()
	query := u.DB.SQ.Update(userTable).
		Set("is_active", true).
		Set("updated_at", now).
		Where(squirrel.Eq{"id": userID})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("UnFreezeUserByID: failed to build query: %w", err)
	}

	commandTag, err := u.DB.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("UnFreezeUserByID: failed to execute query: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("UnFreezeUserByID: user with ID %d not found or already active", userID)
	}

	return nil
}

// FindAllUsers retrieves a paginated list of all users.
func (u *userRepository) FindAllUsers(ctx context.Context, limit, offset uint64) ([]*models.User, error) {
	query := u.DB.SQ.Select("id", "name", "role", "email", "kratos_identity_id", "is_active", "created_at", "updated_at").
		From(userTable).
		OrderBy("name ASC"). // Example ordering
		Limit(limit).
		Offset(offset)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("FindAllUsers: failed to build query: %w", err)
	}

	users := []*models.User{}
	err = pgxscan.Select(ctx, u.DB.Pool, &users, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("FindAllUsers: failed to execute query or scan: %w", err)
	}

	return users, nil
}

// CountUsers counts the total number of users.
func (u *userRepository) CountUsers(ctx context.Context) (int, error) {
	query := u.DB.SQ.Select("COUNT(*)").From(userTable)
	sql, args, err := query.ToSql()
	if err != nil {
		return 0, fmt.Errorf("CountUsers: failed to build query: %w", err)
	}
	var count int
	err = u.DB.Pool.QueryRow(ctx, sql, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("CountUsers: failed to execute query or scan: %w", err)
	}
	return count, nil
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
		Set("updated_at", model.UpdatedAt). // Update timestamp (Removed roll_no)
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
// This implementation updates directly by ID without fetching first.
func (u *userRepository) FreezeUserByID(ctx context.Context, userID int) error {
	// Build the UPDATE query
	now := time.Now()
	query := u.DB.SQ.Update(userTable).
		Set("is_active", false).         // Set status to false
		Set("updated_at", now).          // Update timestamp
		Where(squirrel.Eq{"id": userID}) // Identify the user by ID

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
		return fmt.Errorf("FreezeUserByID: user with ID %d not found or already frozen", userID)
	}

	return nil // Success
}

// DeleteUser deletes a user record based on their roll number.
// This implementation deletes directly by ID without fetching first.
func (u *userRepository) DeleteUserByID(ctx context.Context, userID int) error {
	// Build the DELETE query
	query := u.DB.SQ.Delete(userTable).
		Where(squirrel.Eq{"id": userID}) // Identify the user by ID

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
		return fmt.Errorf("DeleteUserByID: user with ID %d not found", userID)
	}

	return nil // Success
}
