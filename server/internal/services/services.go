package services

import (
	"eduhub/server/internal/config"
	"eduhub/server/internal/repository"
	"eduhub/server/internal/services/attendance"
	"eduhub/server/internal/services/auth"
	"eduhub/server/internal/services/student"
	"eduhub/server/internal/services/system"
)

type Services struct {
	Auth auth.AuthService
	// Quiz *Quiz.QuizService
	// Fee *Fee.FeeService
	Attendance     attendance.AttendanceService
	StudentService student.StudentService
	System         system.SystemService
}

func NewServices(cfg *config.Config) *Services {
	kratosService := auth.NewKratosService()
	ketoService := auth.NewKetoService()
	authService := auth.NewAuthService(kratosService, ketoService)
	db := config.LoadDatabase()
	repository := repository.NewRepository(db)
	studentService := student.NewstudentService(repository.StudentRepository, repository.AttendanceRepository, repository.EnrollmentRepository)
	systemService := system.NewSystemService(db)
	attendanceService := attendance.NewAttendanceService(repository.AttendanceRepository, repository.StudentRepository, repository.EnrollmentRepository)
	return &Services{
		Auth:           authService,
		Attendance:     attendanceService,
		StudentService: studentService,
		System:         systemService,
	}
}
