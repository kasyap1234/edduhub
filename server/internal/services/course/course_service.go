package course

import (
	"context"
	"eduhub/server/internal/models"
	"eduhub/server/internal/repository"
)

type CourseService interface {
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

type courseService struct {
	courseRepo repository.CourseRepository
}

func NewCourseService(courseRepo repository.CourseRepository) CourseService {
	return &courseService{
		courseRepo: courseRepo,
	}
}

func (c *courseService) CreateCourse(ctx context.Context, course *models.Course) error {
	return c.courseRepo.CreateCourse(ctx, course)
}

func (c *courseService) FindCourseByID(ctx context.Context, collegeID int, courseID int) (*models.Course, error) {

	return c.courseRepo.FindCourseByID(ctx, collegeID, courseID)
}
func (c *courseService) UpdateCourse(ctx context.Context, course *models.Course) error {
	return c.courseRepo.UpdateCourse(ctx, course)
}
func (c *courseService) DeleteCourse(ctx context.Context, collegeID int, courseID int) error {
	return c.courseRepo.DeleteCourse(ctx, collegeID, courseID)

}

// Find methods with pagination
func (c *courseService) FindAllCourses(ctx context.Context, collegeID int, limit, offset uint64) ([]*models.Course, error) {
	return c.courseRepo.FindAllCourses(ctx, collegeID, limit, offset)
}
func (c *courseService) FindCoursesByInstructor(ctx context.Context, collegeID int, instructorID int, limit, offset uint64) ([]*models.Course, error) {
	return c.courseRepo.FindCoursesByInstructor(ctx, collegeID, instructorID, limit, offset)
}

// Count methods
func (c *courseService) CountCoursesByCollege(ctx context.Context, collegeID int) (int, error) {
	return c.courseRepo.CountCoursesByCollege(ctx, collegeID)
}
func (c *courseService) CountCoursesByInstructor(ctx context.Context, collegeID int, instructorID int) (int, error) {
	return c.courseRepo.CountCoursesByInstructor(ctx, collegeID, instructorID)
}

