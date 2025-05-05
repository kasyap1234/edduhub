package services

import (
	"eduhub/server/internal/config"
	"eduhub/server/internal/repository"
	"eduhub/server/internal/services/attendance"
	"eduhub/server/internal/services/auth"
	"eduhub/server/internal/services/student"
)

type Services struct {
	Auth auth.AuthService
	// Quiz *Quiz.QuizService
	// Fee *Fee.FeeService
	Attendance     attendance.AttendanceService
	StudentService student.StudentService

	// System         system.SystemService
}

func NewServices(cfg *config.Config) *Services {
	kratosService := auth.NewKratosService()
	ketoService := auth.NewKetoService()
	authService := auth.NewAuthService(kratosService, ketoService)

	repository := repository.NewRepository(cfg.DB)
	studentService := student.NewstudentService(repository.StudentRepository, repository.AttendanceRepository, repository.EnrollmentRepository)
	// systemService := system.NewSystemService(cfg.DB)
	attendanceService := attendance.NewAttendanceService(repository.AttendanceRepository, repository.StudentRepository, repository.EnrollmentRepository)
	return &Services{
		Auth:           authService,
		Attendance:     attendanceService,
		StudentService: studentService,
		// System:         systemService,
	}
}
