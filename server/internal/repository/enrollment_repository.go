package repository

import (
	"context"
	"eduhub/server/internal/models"
	"errors"

	"github.com/uptrace/bun"
)

type EnrollmentRepository interface {
	CreateEnrollment(ctx context.Context, enrollment *models.Enrollment) error
	IsStudentEnrolled(ctx context.Context, studentID, courseID int) error
}

type enrollmentRepository struct {
	db *bun.DB
}

func NewEnrollmentRepository(db *bun.DB) EnrollmentRepository {
	return &enrollmentRepository{
		db: db,
	}
}

func (e *enrollmentRepository) CreateEnrollment(ctx context.Context, enrollment *models.Enrollment) error {
	_, err := e.db.NewInsert().Model(enrollment).Exec(ctx)
	if err != nil {
		errors.New("unable to create enrollment")
		return err
	}
	return nil
}

func (e *enrollmentRepository) IsStudentEnrolled(ctx context.Context, studentID, courseID int) (bool, error) {
	_, err := e.db.NewSelect().
		Model((*models.Enrollment)(nil)).
		// Add the WHERE clause to find the specific link.
		Where("e.student_id = ? AND e.course_id = ?", studentID, courseID).
		Where("student.is_active=?", true).
		Exists(ctx)

	// Handle potential database errors (connection issues, permissions, etc.)
	if err != nil {
		// Note: .Exists() itself handles sql.ErrNoRows internally, returning false, nil in that case.
		// So we only need to worry about other potential DB errors.
		return false, errors.New("error checking enrollment existence")
	}
	return true, nil
}
