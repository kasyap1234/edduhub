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
	// Attendance *AttendanceHandler
}

func NewHandlers(services *services.Services) *Handlers {
	return &Handlers{
		Auth: NewAuthHandler(services.Auth),
	}
}
