package repository

import (
	"context"
	"errors"
	"testing"

	"eduhub/server/internal/models"

	"github.com/Masterminds/squirrel"
	// "github.com/jackc/pgx/v4" // Not needed for current user repo tests, but keep if adding Get methods
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupUserTest(t *testing.T) (pgxmock.PgxPoolIface, *DB, UserRepository, context.Context) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)

	db := &DB{
		Pool: mock,
		SQ:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}

	repo := NewUserRepository(db)
	ctx := context.Background()

	return mock, db, repo, ctx
}

func TestCreateUser(t *testing.T) {
	mock, _, repo, ctx := setupUserTest(t)
	defer mock.Close()

	user := &models.User{
		Name:             "Test User",
		Role:             "student",
		Email:            "test@example.com",
		KratosIdentityID: "kratos-user-id",
		IsActive:         true,
	}
	expectedID := 25

	// Corrected SQL: Removed roll_no, adjusted placeholders
	sqlRegex := `INSERT INTO users \(name,role,email,kratos_identity_id,is_active,created_at,updated_at\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7\) RETURNING id`

	mock.ExpectQuery(sqlRegex). // Removed regexp.QuoteMeta
					WithArgs(
			user.Name,
			user.Role,
			user.Email,
			user.KratosIdentityID,
			user.IsActive,
			// user.RollNo, // Removed
			pgxmock.AnyArg(), // CreatedAt
			pgxmock.AnyArg(), // UpdatedAt
		).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(expectedID))

	err := repo.CreateUser(ctx, user)

	assert.NoError(t, err)
	assert.Equal(t, expectedID, user.ID)
	assert.False(t, user.CreatedAt.IsZero())
	assert.False(t, user.UpdatedAt.IsZero())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateUser_Error(t *testing.T) {
	mock, _, repo, ctx := setupUserTest(t)
	defer mock.Close()

	user := &models.User{Email: "fail@example.com"} // Use other fields for test data
	dbError := errors.New("unique constraint violation")
	// Corrected SQL: Removed roll_no, adjusted placeholders
	sqlRegex := `INSERT INTO users \(name,role,email,kratos_identity_id,is_active,created_at,updated_at\) VALUES \(\$1,\$2,\$3,\$4,\$5,\$6,\$7\) RETURNING id`

	mock.ExpectQuery(sqlRegex). // Removed regexp.QuoteMeta
					WithArgs(user.Name, user.Role, user.Email, user.KratosIdentityID, user.IsActive, pgxmock.AnyArg(), pgxmock.AnyArg()).
					WillReturnError(dbError)

	err := repo.CreateUser(ctx, user)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to execute query or scan ID")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// --- Test UpdateUser ---

func TestUpdateUser(t *testing.T) {
	mock, _, repo, ctx := setupUserTest(t)
	defer mock.Close()

	userToUpdate := &models.User{
		ID:               25,
		Name:             "Updated User Name",
		Role:             "admin",
		Email:            "updated@example.com",
		KratosIdentityID: "kratos-updated-id",
		IsActive:         false,
		// RollNo:           "USER001-UPDATED", // Removed
	}

	// Corrected SQL: Removed roll_no, adjusted placeholders
	sqlRegex := `UPDATE users SET name = \$1, role = \$2, email = \$3, kratos_identity_id = \$4, is_active = \$5, updated_at = \$6 WHERE id = \$7`

	mock.ExpectExec(sqlRegex). // Removed regexp.QuoteMeta
					WithArgs(
			userToUpdate.Name,
			userToUpdate.Role,
			userToUpdate.Email,
			userToUpdate.KratosIdentityID,
			userToUpdate.IsActive,
			// userToUpdate.RollNo, // Removed
			pgxmock.AnyArg(), // UpdatedAt
			userToUpdate.ID,
		).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err := repo.UpdateUser(ctx, userToUpdate)

	assert.NoError(t, err)
	assert.False(t, userToUpdate.UpdatedAt.IsZero())
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestUpdateUser_NotFound(t *testing.T) {
	mock, _, repo, ctx := setupUserTest(t)
	defer mock.Close()

	userToUpdate := &models.User{ID: 999}
	// Corrected SQL: Removed roll_no, adjusted placeholders
	sqlRegex := `UPDATE users SET name = \$1, role = \$2, email = \$3, kratos_identity_id = \$4, is_active = \$5, updated_at = \$6 WHERE id = \$7`

	mock.ExpectExec(sqlRegex). // Removed regexp.QuoteMeta
					WithArgs(userToUpdate.Name, userToUpdate.Role, userToUpdate.Email, userToUpdate.KratosIdentityID, userToUpdate.IsActive, pgxmock.AnyArg(), userToUpdate.ID).
					WillReturnResult(pgxmock.NewResult("UPDATE", 0))

	err := repo.UpdateUser(ctx, userToUpdate)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no row updated")
	assert.NoError(t, mock.ExpectationsWereMet())
}

// --- Test FreezeUserByID ---

func TestFreezeUserByID(t *testing.T) {
	mock, _, repo, ctx := setupUserTest(t)
	defer mock.Close()

	userID := 25
	// Corrected SQL: WHERE clause uses ID
	sqlRegex := `UPDATE users SET is_active = \$1, updated_at = \$2 WHERE id = \$3`

	mock.ExpectExec(sqlRegex). // Removed regexp.QuoteMeta
					WithArgs(false, pgxmock.AnyArg(), userID).
					WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	err := repo.FreezeUserByID(ctx, userID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFreezeUserByID_NotFound(t *testing.T) {
	mock, _, repo, ctx := setupUserTest(t)
	defer mock.Close()

	userID := 999
	// Corrected SQL: WHERE clause uses ID
	sqlRegex := `UPDATE users SET is_active = \$1, updated_at = \$2 WHERE id = \$3`

	mock.ExpectExec(sqlRegex). // Removed regexp.QuoteMeta
					WithArgs(false, pgxmock.AnyArg(), userID).
					WillReturnResult(pgxmock.NewResult("UPDATE", 0))

	err := repo.FreezeUserByID(ctx, userID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found or already frozen") // Error message from repo function
	assert.NoError(t, mock.ExpectationsWereMet())
}

// --- Test DeleteUserByID --- // Renamed test function

func TestDeleteUserByID(t *testing.T) { // Renamed test function
	mock, _, repo, ctx := setupUserTest(t)
	defer mock.Close()

	userID := 25
	// Corrected SQL: WHERE clause uses ID
	sqlRegex := `DELETE FROM users WHERE id = \$1`

	mock.ExpectExec(sqlRegex). // Removed regexp.QuoteMeta
					WithArgs(userID).
					WillReturnResult(pgxmock.NewResult("DELETE", 1))

	err := repo.DeleteUserByID(ctx, userID)

	assert.NoError(t, err)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteUserByID_NotFound(t *testing.T) { // Renamed test function
	mock, _, repo, ctx := setupUserTest(t)
	defer mock.Close()

	userID := 999
	// Corrected SQL: WHERE clause uses ID
	sqlRegex := `DELETE FROM users WHERE id = \$1`

	mock.ExpectExec(sqlRegex). // Removed regexp.QuoteMeta
					WithArgs(userID).
					WillReturnResult(pgxmock.NewResult("DELETE", 0))

	err := repo.DeleteUserByID(ctx, userID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestDeleteUserByID_Error(t *testing.T) { // Renamed test function
	mock, _, repo, ctx := setupUserTest(t)
	defer mock.Close()

	userID := 25
	dbError := errors.New("delete constraint error")
	// Corrected SQL: WHERE clause uses ID
	sqlRegex := `DELETE FROM users WHERE id = \$1`

	mock.ExpectExec(sqlRegex). // Removed regexp.QuoteMeta
					WithArgs(userID).
					WillReturnError(dbError)

	err := repo.DeleteUserByID(ctx, userID)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to execute query")
	assert.NoError(t, mock.ExpectationsWereMet())
}
