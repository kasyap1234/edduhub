package handler

import (
	"eduhub/server/internal/services"
)

type Handlers struct {
	Auth *AuthHandler
	// other handlers
	// quiz handler
	// fee handler
	// attendance handler
	Attendance *AttendanceHandler
}

func NewHandlers(services *services.Services) *Handlers {
	return &Handlers{

		Auth:       NewAuthHandler(services.Auth, services.StudentService),
		Attendance: NewAttedanceHandler(services.Attendance),
		// other handlers

	}
}
