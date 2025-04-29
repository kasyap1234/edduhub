package handler

import (
	"eduhub/server/internal/services"
	// Import system service package if needed for type
)

type Handlers struct {
	Auth *AuthHandler
	// other handlers
	// quiz handler
	// fee handler
	// attendance handler
	Attendance *AttendanceHandler
	// System     *SystemHandler
}

func NewHandlers(services *services.Services) *Handlers {
	return &Handlers{
		Auth:       NewAuthHandler(services.Auth),
		Attendance: NewAttendanceHandler(services.Attendance),
		// other handlers
		// System: NewSystemHandler(services.System),
	}
}
