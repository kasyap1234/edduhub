package enrollment

import (
	"context"

	"eduhub/server/internal/models"
)

type EnrollmentService interface {
	CreateEnrollment(ctx context.Context, enrollment *models.Enrollment) error
	IsStudentEnrolled(ctx context.Context, collegeID, studentID, courseID int) (bool, error)
	UpdateEnrollment(ctx context.Context, enrollment *models.Enrollment) error
	UpdateEnrollmentStatus(ctx context.Context, id int, status string) error
	DeleteEnrollment(ctx context.Context, collegeID int, enrollmentID int) error
	FindEnrollmentsByStudent(ctx context.Context, collegeID int, studentID int) ([]*models.Enrollment, error)
}
