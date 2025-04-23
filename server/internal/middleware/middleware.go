package middleware

import (
	"eduhub/server/internal/repository"
	"eduhub/server/internal/services"

)

type Middleware struct {
	Auth *AuthMiddleware
	// other middleware
}

func NewMiddleware(services *services.Services) *Middleware {
	authSvc := services.Auth
	studentRepo := repos.StudentRepository
	enrollmentRepo := repos.EnrollmentRepository
	attendanceRepo := repos.AttendanceRepository
	studentService :=services.studentService
	return &Middleware{

		Auth: NewAuthMiddleware(authSvc, studentService),
	}

}
