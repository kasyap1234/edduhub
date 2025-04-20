package student

import (
	"context"
	"eduhub/server/internal/models"
	"eduhub/server/internal/repository"
)

type StudentService interface {
	CreateStudent(ctx context.Context, student *models.Student) error
	GetStudentByRollNo(ctx context.Context, rollNo string) (*models.Student, error)
	GetStudentByID(ctx context.Context, studentID int) (*models.Student, error)
	UpdateStudent(ctx context.Context, model *models.Student) error
	FreezeStudent(ctx context.Context, RollNo string) error
	UnFreezeStudent(ctx context.Context, RollNo string) error
	FindByKratosID(ctx context.Context, kratosID string) (*models.Student, error)
	VerifyStudentEnrollment(ctx context.Context, collegeID, studentID, courseID int) (bool, error)
}

type studentService struct {
	StudentRepo    repository.StudentRepository
	AttendanceRepo repository.AttendanceRepository
}

func NewstudentService(studentRepo repository.StudentRepository) StudentService {
	return &studentService{
		StudentRepo: studentRepo,
	}
}
func (s *studentService) CreateStudent(ctx context.Context, student *models.Student) error {
	return s.StudentRepo.CreateStudent(ctx, student)
}

func (s *studentService) FindByKratosID(ctx context.Context, kratosID string) (*models.Student, error) {
	student, err := s.StudentRepo.FindByKratosID(ctx, kratosID)
	return student, err
}

func (s *studentService) GetStudentByRollNo(ctx context.Context, rollNo string) (*models.Student, error) {
	student, err := s.StudentRepo.GetStudentByRollNo(ctx, rollNo)
	return student, err
}

func (s *studentService) GetStudentByID(ctx context.Context, studentID int) (*models.Student, error) {
	student, err := s.StudentRepo.GetStudentByID(ctx, studentID)
	return student, err
}

func (s *studentService) UpdateStudent(ctx context.Context, model *models.Student) error {
	return s.StudentRepo.UpdateStudent(ctx, model)
}

func (s *studentService) FreezeStudent(ctx context.Context, RollNo string) error {
	return s.StudentRepo.FreezeStudent(ctx, RollNo)
}

func (s *studentService) UnFreezeStudent(ctx context.Context, RollNo string) error {
	return s.StudentRepo.UnFreezeStudent(ctx, RollNo)
}

func (s *studentService) VerifyStudentEnrollment(ctx context.Context, collegeID, studentID, courseID int) (bool, error) {
	return s.StudentRepo.VerifyStudentEnrollment(ctx, collegeID, studentID, courseID)
}
