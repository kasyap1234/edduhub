package repository

import (
	"context"
	"errors"
	"fmt" // Keep fmt for error wrapping
	"time"

	"eduhub/server/internal/models"

	"github.com/Masterminds/squirrel"
	"github.com/georgysavva/scany/pgxscan" // Add pgxscan import
	"github.com/jackc/pgx/v4"              // Use v4 for pgx.ErrNoRows
)

type CourseRepository interface {
	CreateCourse(ctx context.Context, course *models.Course) error
	FindCourseByID(ctx context.Context, collegeID int, courseID int) (*models.Course, error) // Added collegeID
	UpdateCourse(ctx context.Context, course *models.Course) error
	DeleteCourse(ctx context.Context, collegeID int, courseID int) error

	// Find methods with pagination
	FindAllCourses(ctx context.Context, collegeID int, limit, offset uint64) ([]*models.Course, error)
	FindCoursesByInstructor(ctx context.Context, collegeID int, instructorID int, limit, offset uint64) ([]*models.Course, error)

	// Count methods
	CountCoursesByCollege(ctx context.Context, collegeID int) (int, error)
	CountCoursesByInstructor(ctx context.Context, collegeID int, instructorID int) (int, error)
}

type courseRepository struct {
	DB *DB
}

func NewCourseRepository(db *DB) CourseRepository {
	return &courseRepository{
		DB: db,
	}
}

const courseTable = "course" // Define table name constant

func (c *courseRepository) CreateCourse(ctx context.Context, course *models.Course) error {
	// Set timestamps
	now := time.Now()
	if course.CreatedAt.IsZero() {
		course.CreatedAt = now
	}
	if course.UpdatedAt.IsZero() {
		course.UpdatedAt = now
	}

	query := c.DB.SQ.Insert(courseTable).
		Columns("name", "description", "credits", "instructor_id", "college_id", "created_at", "updated_at"). // Added college_id, created_at, updated_at
		Values(
			course.Name,
			course.Description,
			course.Credits,
			course.InstructorID,
			course.CollegeID, // Assuming CollegeID exists in models.Course
			course.CreatedAt,
			course.UpdatedAt,
		).
		Suffix("RETURNING id")
	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("create course query build error: %w", err)
	}

	err = c.DB.Pool.QueryRow(ctx, sql, args...).Scan(&course.ID) // Scan the returned ID into the original course struct
	if err != nil {
		// Consider checking for specific DB errors (e.g., unique constraints)
		return fmt.Errorf("CreateCourse: failed to execute query or scan ID: %w", err)
	}
	return nil
}

func (c *courseRepository) FindCourseByID(ctx context.Context, collegeID int, courseID int) (*models.Course, error) {
	query := c.DB.SQ.Select(
		"id", "name", "description", "credits", "instructor_id", "college_id", "created_at", "updated_at", // Added college_id, timestamps
	).
		From(courseTable).
		Where(squirrel.Eq{
			"id":         courseID,
			"college_id": collegeID, // Scope by college
		})

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("FindCourseByID: failed to build query: %w", err)
	}

	// Initialize the struct BEFORE scanning!
	course := &models.Course{}

	// Use pgxscan.Get for single row
	err = pgxscan.Get(ctx, c.DB.Pool, course, sql, args...)
	if err != nil {
		// It's better to check for specific errors like "no rows"
		if errors.Is(err, pgx.ErrNoRows) { // Use errors.Is for checking pgx.ErrNoRows
			return nil, fmt.Errorf("FindCourseByID: course with ID %d not found for college ID %d", courseID, collegeID) // Or a custom ErrNotFound
		}
		return nil, fmt.Errorf("unable to find course: %w", err) // Wrap the original error
	}

	return course, nil
}

func (c *courseRepository) UpdateCourse(ctx context.Context, course *models.Course) error {
	course.UpdatedAt = time.Now()

	query := c.DB.SQ.Update(courseTable).
		Set("name", course.Name).
		Set("description", course.Description).
		Set("credits", course.Credits).
		Set("instructor_id", course.InstructorID).
		Set("updated_at", course.UpdatedAt).
		Where(squirrel.Eq{
			"id":         course.ID,
			"college_id": course.CollegeID, // Ensure update is scoped
		})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("UpdateCourse: failed to build query: %w", err)
	}

	commandTag, err := c.DB.Pool.Exec(ctx, sql, args...)
	if err != nil {
		return fmt.Errorf("UpdateCourse: failed to execute query: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("UpdateCourse: no course found with ID %d for college ID %d, or no changes made", course.ID, course.CollegeID)
	}

	return nil
}

func (c *courseRepository) DeleteCourse(ctx context.Context, collegeID int, courseID int) error {
	query := c.DB.SQ.Delete(courseTable).
		Where(squirrel.Eq{
			"id":         courseID,
			"college_id": collegeID, // Scope deletion
		})

	sql, args, err := query.ToSql()
	if err != nil {
		return fmt.Errorf("DeleteCourse: failed to build query: %w", err)
	}

	commandTag, err := c.DB.Pool.Exec(ctx, sql, args...)
	if err != nil {
		// Consider foreign key constraint errors
		return fmt.Errorf("DeleteCourse: failed to execute query: %w", err)
	}

	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("DeleteCourse: no course found with ID %d for college ID %d, or already deleted", courseID, collegeID)
	}

	return nil
}

func (c *courseRepository) FindAllCourses(ctx context.Context, collegeID int, limit, offset uint64) ([]*models.Course, error) {
	query := c.DB.SQ.Select(
		"id", "name", "description", "credits", "instructor_id", "college_id", "created_at", "updated_at",
	).
		From(courseTable).
		Where(squirrel.Eq{"college_id": collegeID}).
		OrderBy("name ASC"). // Example ordering
		Limit(limit).
		Offset(offset)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("FindAllCourses: failed to build query: %w", err)
	}

	courses := []*models.Course{}
	err = pgxscan.Select(ctx, c.DB.Pool, &courses, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("FindAllCourses: failed to execute query or scan: %w", err)
	}

	return courses, nil
}

func (c *courseRepository) FindCoursesByInstructor(ctx context.Context, collegeID int, instructorID int, limit, offset uint64) ([]*models.Course, error) {
	query := c.DB.SQ.Select(
		"id", "name", "description", "credits", "instructor_id", "college_id", "created_at", "updated_at",
	).
		From(courseTable).
		Where(squirrel.Eq{
			"college_id":    collegeID,
			"instructor_id": instructorID,
		}).
		OrderBy("name ASC").
		Limit(limit).
		Offset(offset)

	sql, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("FindCoursesByInstructor: failed to build query: %w", err)
	}

	courses := []*models.Course{}
	err = pgxscan.Select(ctx, c.DB.Pool, &courses, sql, args...)
	if err != nil {
		return nil, fmt.Errorf("FindCoursesByInstructor: failed to execute query or scan: %w", err)
	}

	return courses, nil
}

func (c *courseRepository) CountCoursesByCollege(ctx context.Context, collegeID int) (int, error) {
	return c.countCourses(ctx, squirrel.Eq{"college_id": collegeID})
}

func (c *courseRepository) CountCoursesByInstructor(ctx context.Context, collegeID int, instructorID int) (int, error) {
	return c.countCourses(ctx, squirrel.Eq{"college_id": collegeID, "instructor_id": instructorID})
}

// countCourses is a helper function for counting based on conditions.
func (c *courseRepository) countCourses(ctx context.Context, whereClause squirrel.Sqlizer) (int, error) {
	query := c.DB.SQ.Select("COUNT(*)").From(courseTable).Where(whereClause)
	sql, args, err := query.ToSql()
	if err != nil {
		return 0, fmt.Errorf("countCourses: failed to build query: %w", err)
	}
	var count int
	err = c.DB.Pool.QueryRow(ctx, sql, args...).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("countCourses: failed to execute query or scan: %w", err)
	}
	return count, nil
}
