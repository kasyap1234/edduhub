package repository

import (
	"context"
	"eduhub/server/internal/models"
)

type EnrollmentRepository interface {
	CreateEnrollment(ctx context.Context, enrollment *models.Enrollment) error
	IsStudentEnrolled(ctx context.Context, collegeID, studentID, courseID int) (bool, error)
}

type enrollmentRepository struct {
	db DatabaseRepository[models.Enrollment]
}

func NewEnrollmentRepository(db DatabaseRepository[models.Enrollment]) EnrollmentRepository {
	return &enrollmentRepository{
		db: db,
	}
}

func (e *enrollmentRepository) CreateEnrollment(ctx context.Context, enrollment *models.Enrollment) error {
	err := e.db.Create(ctx, enrollment)
	return err
}

func (e *enrollmentRepository) IsStudentEnrolled(ctx context.Context, collegeID, studentID, courseID int) (bool, error) {

	exists, err := e.db.Exists(ctx, (*models.Enrollment)(nil), "college_id=? AND student_id=? AND course_id=?", collegeID, studentID, courseID)
	return exists, err
}
