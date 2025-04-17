package middleware

import (
	"eduhub/server/internal/services"
	"eduhub/server/internal/services/auth"
)

type Middleware struct {
	authmiddleware *AuthMiddleware
	// other middleware
}

func NewMiddleware(services *services.Services) *Middleware {
	return &Middleware{
		authmiddleware: NewAuthMiddleware(auth.NewAuthService(auth.NewKratosService(), auth.NewKetoService())),
	}

}
