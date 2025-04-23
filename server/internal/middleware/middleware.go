package middleware

import (
	"eduhub/server/internal/repository"
	"eduhub/server/internal/services"
	"eduhub/server/internal/services/auth"
	"eduhub/server/internal/services/student"
)

type Middleware struct {
	Auth *AuthMiddleware
	// other middleware
}

func NewMiddleware(services *services.Services, repos *repository.Repository) *Middleware {
	authSvc := auth.NewAuthService(auth.NewKratosService(), auth.NewKetoService())
	studentRepo := repos.StudentRepository
	enrollmentRepo := repos.EnrollmentRepository
	attendanceRepo := repos.AttendanceRepository
	studentService := student.NewstudentService(studentRepo, attendanceRepo, enrollmentRepo)
	return &Middleware{

		Auth: NewAuthMiddleware(authSvc, studentService),
	}

}
