package repository

import (
	"context"
	"eduhub/server/internal/models"
)

type CourseRepository interface {
	CreateCourse(ctx context.Context, course *models.Course) error
}

type courseRepository struct {
	db DatabaseRepository[models.Course]
}

func NewCourseRepository(db DatabaseRepository[models.Course]) CourseRepository {
	return &courseRepository{
		db: db,
	}
}

func (c *courseRepository) CreateCourse(ctx context.Context, course *models.Course) error {
	return c.db.Create(ctx, course)
}

func (c *courseRepository) FindCourseByID(ctx context.Context, courseID int) (*models.Course, error) {
	return c.db.FindByID(ctx, courseID)
}
