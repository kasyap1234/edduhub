package repository

import (
	"context"
	"errors"
	"fmt"

	"eduhub/server/internal/models"

	"github.com/jackc/pgx"
	"github.com/georgysavva/scany/pgxscan" // Add pgxscan import
)

type CourseRepository interface {
	CreateCourse(ctx context.Context, course *models.Course) error
	FindCourseByID(ctx context.Context, courseID int) (*models.Course, error)
	UpdateCourse(ctx context.Context,course *models.Course)
}

type courseRepository struct {
	DB *DB
}

func NewCourseRepository(db *DB) CourseRepository {
	return &courseRepository{
		DB: db,
	}
}

func (c *courseRepository) CreateCourse(ctx context.Context, course *models.Course) error {
	// Remove "ID" from Columns and course.ID from Values
	// Change Suffix to "RETURNING id" to get only the generated ID
	query := c.DB.SQ.Insert("course").
		Columns("Name", "Description", "Credits", "InstructorID").
		Values(course.Name, course.Description, course.Credits, course.InstructorID).
		Suffix("RETURNING id")
	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("create course query build error: %w", err)
	}

	err = c.DB.Pool.QueryRow(ctx, sql, args...).Scan(&course.ID) // Scan the returned ID into the original course struct
	if err != nil {
		return errors.New("unable to create a course")
	}
	return nil
}

func (c *courseRepository) FindCourseByID(ctx context.Context, courseID int) (*models.Course, error) {
	// Ensure your model field names match the column names or use aliases
	// if they are different and you were using a scanning helper library.
	// With pgxpool.QueryRow/Scan directly, the order and type must match.
	query := c.DB.SQ.Select("ID", "Name", "Description", "Credits", "InstructorID"). // Select "ID" column
												From("course").
												Where("ID = ?", courseID) // Or use squirrel.Eq{"ID": courseID}

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("find course by id query build error: %w", err)
	}

	// Initialize the struct BEFORE scanning!
	course := &models.Course{}

	// Use pgxscan.Get for consistency
	err = pgxscan.Get(ctx, c.DB.Pool, course, sql, args...)
	if err != nil {
		// It's better to check for specific errors like "no rows"
		if err == pgx.ErrNoRows { // Make sure you've imported "github.com/jackc/pgx/v4"
			return nil, fmt.Errorf("course with ID %d not found", courseID) // Or a custom ErrNotFound
		}
		return nil, fmt.Errorf("unable to find course: %w", err) // Wrap the original error
	}

	return course, nil
}
