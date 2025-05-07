package student

import (
	"context"

	"eduhub/server/internal/models"
	"eduhub/server/internal/repository"
)

type StudentService interface {
	FindByKratosID(ctx context.Context, kratosID string) (*models.Student, error)
}

type studentService struct {
	studentRepo    repository.StudentRepository
	attendanceRepo repository.AttendanceRepository
	enrollmentRepo repository.EnrollmentRepository
}

func NewstudentService(studentRepo repository.StudentRepository, attendanceRepo repository.AttendanceRepository, enrollmentRepo repository.EnrollmentRepository) StudentService {
	return &studentService{
		studentRepo:    studentRepo,
		attendanceRepo: attendanceRepo,
		enrollmentRepo: enrollmentRepo,
	}
}

func (a *studentService) FindByKratosID(ctx context.Context, kratosID string) (*models.Student, error) {
	student, err := a.studentRepo.FindByKratosID(ctx, kratosID)
	if err != nil {
		return nil, err
	}
	return student, nil
}
