package handler

import (
	"eduhub/server/internal/middleware"

	"github.com/labstack/echo/v4"
)

func SetupRoutes(e *echo.Echo, a *Handlers, m *middleware.AuthMiddleware) {

	// auth routes

	// protected college routes
	// protected finance routes
	auth := e.Group("/auth")
	auth.POST("/register", a.Auth.InitiateRegistration)
	auth.POST("/auth/register/complete", a.Auth.HandleRegistration)
	auth.GET("/login", a.Auth.HandleLogin)
	auth.GET("/callback", a.Auth.HandleCallback)
	
	apiGroup := e.Group("/api", m.ValidateSession, m.RequireCollege)

	attendance := apiGroup.Group("/attendance")
	attendance.POST("/", a.Attendance.MarkAttendance, m.RequireRole(middleware.RoleAdmin, middleware.RoleFaculty, middleware.RoleStudent))
	attendance.GET("/get-attendance-course", a.Attendance.GetAttendanceByCourse, m.RequireRole(middleware.RoleAdmin, middleware.RoleFaculty))
	attendance.GET("/student/:studentID", a.Attendance.GetAttendanceForStudent, m.RequireRole(middleware.RoleAdmin, middleware.RoleFaculty, middleware.RoleStudent))
	attendance.GET("/student/:studentID/course/:courseID", a.Attendance.GetAttendanceByStudentAndCourse, m.RequireRole(middleware.RoleAdmin, middleware.RoleFaculty, middleware.RoleStudent))
	attendance.GET("course", a.Attendance.GetAttendanceByCourse, m.RequireRole(middleware.RoleAdmin, middleware.RoleFaculty))
	// to do add multitenancy properly
}
