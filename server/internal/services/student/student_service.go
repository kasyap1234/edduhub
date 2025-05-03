package student

import "eduhub/server/internal/repository"

type StudentService interface {
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
