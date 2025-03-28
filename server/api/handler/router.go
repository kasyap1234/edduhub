package handler

import (
	"eduhub/server/internal/middleware"
	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, a *Handlers, m *middleware.Middleware) {

	// auth routes

	// protected college routes
	// protected finance routes
	auth := e.Group("/auth")
	auth.POST("/register", a.Auth.InitiateRegistration)
	auth.POST("/auth/register/complete", a.Auth.HandleRegistration)
	auth.GET("/login", a.Auth.HandleLogin)
	auth.GET("/callback", a.Auth.HandleCallback)

	attendance := e.Group("/attendance")
	attendance.POST("/mark", a.Attendance.MarkAttendance)
	attendance.GET("/student/{id}")
	attendance.GET("/course/{id}")
}
