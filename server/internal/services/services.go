package services

import (
	"eduhub/server/internal/config"
	"eduhub/server/internal/repository"
	"eduhub/server/internal/services/attendance"
	"eduhub/server/internal/services/auth"
)

type Services struct {
	Auth auth.AuthService
	// Quiz *Quiz.QuizService
	// Fee *Fee.FeeService
	Attendance attendance.AttendanceService
}

func NewServices(cfg *config.Config) *Services {
	kratosService := auth.NewKratosService()
	ketoService := auth.NewKetoService()
	authService := auth.NewAuthService(kratosService, ketoService)
	db := config.LoadDatabase()
	repository := repository.NewRepository(db)
	attendanceService := attendance.NewAttendanceService(repository.AttendanceRepository, repository.StudentRepository)
	return &Services{
		Auth:       authService,
		Attendance: attendanceService,
	}

}
