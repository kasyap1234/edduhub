package repository_test

import (
	"context"
	"database/sql"
	"eduhub/server/internal/models"
	"eduhub/server/internal/repository"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

func setupTestDB(t *testing.T) *bun.DB {
	// Connect to the test database
	dsn := "postgres://user:password@localhost:5432/testdb?sslmode=disable"
	sqldb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(dsn)))
	db := bun.NewDB(sqldb, pgdialect.New())

	// Create the students table
	_, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS students (
            student_id SERIAL PRIMARY KEY,
            kratos_identity_id TEXT NOT NULL,
            college_id INTEGER,
            roll_no TEXT,
            batch INTEGER,
            year INTEGER,
            sem INTEGER,
            is_active BOOLEAN
        )
    `)
	assert.NoError(t, err)

	// Truncate the table before each test
	_, err = db.Exec("TRUNCATE TABLE students RESTART IDENTITY")
	assert.NoError(t, err)

	return db
}

func TestFindByKratosID(t *testing.T) {
	// Arrange
	ctx := context.Background()
	db := setupTestDB(t)
	defer db.Close()

	studentRepo := repository.NewBaseRepository[models.Student](db)
	repo := repository.NewStudentRepository(studentRepo)
	// Insert test data
	_, err := db.NewInsert().Model(&models.Student{
		KratosIdentityID: "123",
		CollegeID:        3,
		RollNo:           "SE20UARI73",
		Batch:            2020,
		Year:             5,
		Sem:              2,
		IsActive:         true,
	}).Exec(ctx)
	assert.NoError(t, err)

	// Define test cases
	testCases := []struct {
		name           string
		kratosID       string
		expectedResult *models.Student
		expectedError  bool
	}{
		{
			name:     "Success case",
			kratosID: "123",
			expectedResult: &models.Student{
				KratosIdentityID: "123",
				CollegeID:        3,
				RollNo:           "SE20UARI73",
				Batch:            2020,
				Year:             5,
				Sem:              2,
				IsActive:         true,
			},
			expectedError: false,
		},
		{
			name:           "Failure case - Not found",
			kratosID:       "999",
			expectedResult: nil,
			expectedError:  true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			result, err := repo.FindByKratosID(ctx, tc.kratosID)

			// Assert
			if tc.expectedError {
				assert.Error(t, err, "Expected an error but got none")
			} else {
				assert.NoError(t, err, "Expected no error but got one")
				assert.Equal(t, tc.expectedResult.KratosIdentityID, result.KratosIdentityID, "Expected result does not match actual result")
			}
		})
	}
}

func TestCreateStudent(t *testing.T) {
	// Arrange
	ctx := context.Background()
	db := setupTestDB(t)
	defer db.Close()

	studentRepo := repository.NewBaseRepository[models.Student](db)
	repo := repository.NewStudentRepository(studentRepo)
	// Define test cases
	testCases := []struct {
		name          string
		student       *models.Student
		expectedError bool
	}{
		{
			name: "Success case",
			student: &models.Student{
				KratosIdentityID: "identity",
				CollegeID:        3,
				RollNo:           "SE20UARI73",
				Batch:            2020,
				Year:             5,
				Sem:              2,
				IsActive:         true,
			},
			expectedError: false,
		},
		{
			name: "Failure case - Duplicate entry",
			student: &models.Student{
				KratosIdentityID: "identity", // Duplicate ID
			},
			expectedError: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Act
			err := repo.CreateStudent(ctx, tc.student)

			// Assert
			if tc.expectedError {
				assert.Error(t, err, "Expected an error but got none")
			} else {
				assert.NoError(t, err, "Expected no error but got one")
			}
		})
	}
}
