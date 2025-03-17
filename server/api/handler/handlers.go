package handler

import (
	"eduhub/server/internal/services"
	"eduhub/server/internal/services/auth"
)

type Handlers struct {
	Auth *AuthHandler
	// other handlers
	// quiz handler
	// fee handler
	// attendance handler
}

func NewHandlers(services *services.Services) *Handlers {
	return &Handlers{
		Auth: NewAuthHandler(auth.NewKratosService()),
	}
}
