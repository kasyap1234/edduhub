package repository

import (
	"context"
	"errors"
	"regexp" // Import regexp for ExpectQuery with regex
	"testing"

	"eduhub/server/internal/models"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4"
	"github.com/pashagolub/pgxmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupCourseTest(t *testing.T) (pgxmock.PgxPoolIface, *DB, CourseRepository, context.Context) {
	mock, err := pgxmock.NewPool()
	require.NoError(t, err)

	db := &DB{
		Pool: mock,
		SQ:   squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}

	repo := NewCourseRepository(db)
	ctx := context.Background()

	return mock, db, repo, ctx
}

func TestCreateCourse(t *testing.T) {
	mock, _, repo, ctx := setupCourseTest(t)
	defer mock.Close()

	course := &models.Course{
		Name:         "Introduction to Testing",
		Description:  "Learn how to test Go code",
		Credits:      3,
		InstructorID: 1,
	}
	expectedID := 10

	// Use regex to match the SQL query, escaping special characters
	// This allows flexibility if the exact spacing or casing changes slightly
	sqlRegex := `INSERT INTO course \(Name,Description,Credits,InstructorID\) VALUES \(\$1,\$2,\$3,\$4\) RETURNING id`

	mock.ExpectQuery(regexp.QuoteMeta(sqlRegex)). // Use regexp.QuoteMeta if matching exact string, or provide regex pattern
							WithArgs(course.Name, course.Description, course.Credits, course.InstructorID).
							WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(expectedID))

	err := repo.CreateCourse(ctx, course)

	assert.NoError(t, err)
	assert.Equal(t, expectedID, course.ID)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestCreateCourse_Error(t *testing.T) {
	mock, _, repo, ctx := setupCourseTest(t)
	defer mock.Close()

	course := &models.Course{
		Name:         "Error Course",
		Description:  "This will fail",
		Credits:      1,
		InstructorID: 1,
	}
	dbError := errors.New("database error")

	sqlRegex := `INSERT INTO course \(Name,Description,Credits,InstructorID\) VALUES \(\$1,\$2,\$3,\$4\) RETURNING id`

	mock.ExpectQuery(regexp.QuoteMeta(sqlRegex)).
		WithArgs(course.Name, course.Description, course.Credits, course.InstructorID).
		WillReturnError(dbError)

	err := repo.CreateCourse(ctx, course)

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "unable to create a course") // Check against the specific error returned by the repo
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindCourseByID(t *testing.T) {
	mock, _, repo, ctx := setupCourseTest(t)
	defer mock.Close()

	courseID := 10
	expectedCourse := &models.Course{
		ID:           courseID,
		Name:         "Found Course",
		Description:  "Successfully retrieved",
		Credits:      4,
		InstructorID: 2,
	}

	// Use regex for flexibility or exact string match
	sqlRegex := `SELECT ID, Name, Description, Credits, InstructorID FROM course WHERE ID = \$1`
	rows := pgxmock.NewRows([]string{"ID", "Name", "Description", "Credits", "InstructorID"}).
		AddRow(expectedCourse.ID, expectedCourse.Name, expectedCourse.Description, expectedCourse.Credits, expectedCourse.InstructorID)

	mock.ExpectQuery(regexp.QuoteMeta(sqlRegex)).
		WithArgs(courseID).
		WillReturnRows(rows)

	course, err := repo.FindCourseByID(ctx, courseID)

	assert.NoError(t, err)
	assert.Equal(t, expectedCourse, course)
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindCourseByID_NotFound(t *testing.T) {
	mock, _, repo, ctx := setupCourseTest(t)
	defer mock.Close()

	courseID := 99
	sqlRegex := `SELECT ID, Name, Description, Credits, InstructorID FROM course WHERE ID = \$1`

	mock.ExpectQuery(regexp.QuoteMeta(sqlRegex)).
		WithArgs(courseID).
		WillReturnError(pgx.ErrNoRows)

	course, err := repo.FindCourseByID(ctx, courseID)

	assert.Error(t, err)
	assert.Nil(t, course)
	assert.Contains(t, err.Error(), "not found") // Check the specific error message
	assert.NoError(t, mock.ExpectationsWereMet())
}

func TestFindCourseByID_Error(t *testing.T) {
	mock, _, repo, ctx := setupCourseTest(t)
	defer mock.Close()

	courseID := 10
	dbError := errors.New("database connection lost")
	sqlRegex := `SELECT ID, Name, Description, Credits, InstructorID FROM course WHERE ID = \$1`

	mock.ExpectQuery(regexp.QuoteMeta(sqlRegex)).
		WithArgs(courseID).
		WillReturnError(dbError)

	course, err := repo.FindCourseByID(ctx, courseID)

	assert.Error(t, err)
	assert.Nil(t, course)
	assert.Contains(t, err.Error(), "unable to find course") // Check the specific error message
	assert.NoError(t, mock.ExpectationsWereMet())
}