package middleware

import (
	"eduhub/server/internal/repository"
	"eduhub/server/internal/services"
	"eduhub/server/internal/services/auth"
)

type Middleware struct {
	Auth *AuthMiddleware
	// other middleware
}

func NewMiddleware(services *services.Services, repos *repository.Repository) *Middleware {
	authSvc := auth.NewAuthService(auth.NewKratosService(), auth.NewKetoService())
	studentRepo := repos.StudentRepository
	return &Middleware{

		Auth: NewAuthMiddleware(authSvc, studentRepo),
	}

}
