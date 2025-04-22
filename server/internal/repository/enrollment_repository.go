package repository

import (
	"context"
	"eduhub/server/internal/models"
)

type EnrollmentRepository interface {
	CreateEnrollment(ctx context.Context, enrollment *models.Enrollment) error
	IsStudentEnrolled(ctx context.Context, studentID, courseID int) (bool, error)
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

func (e *enrollmentRepository) IsStudentEnrolled(ctx context.Context, studentID, courseID int) (bool, error) {
	// _, err := e.db.NewSelect().
	// 	Model((*models.Enrollment)(nil)).
	// 	// Add the WHERE clause to find the specific link.
	// 	Where("e.student_id = ? AND e.course_id = ?", studentID, courseID).
	// 	Where("student.is_active=?", true).
	// 	Exists(ctx)

	// // Handle potential database errors (connection issues, permissions, etc.)
	// if err != nil {
	// 	// Note: .Exists() itself handles sql.ErrNoRows internally, returning false, nil in that case.
	// 	// So we only need to worry about other potential DB errors.
	// 	return false, errors.New("error checking enrollment existence")
	// }
	// return true, nil
	exists, err := e.db.Exists(ctx, (*models.Enrollment)(nil), "student_id=? AND course_id=?", studentID, courseID)
	return exists, err
}
