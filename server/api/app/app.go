package app

import (
	"eduhub/server/api/handler"
	"eduhub/server/internal/config"
	"eduhub/server/internal/services/auth"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/uptrace/bun"
	"eduhub/server/internal/middleware"
)

type App struct {
	e        *echo.Echo
	db       *bun.DB
	config   *config.Config
	handlers *handler.Handlers
}

func New() *App {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	// Initialize auth service
	authService := auth.NewKratosService()


	// Initialize handlers
	authHandler:=handlers.

	return &App{
		e:        echo.New(),
		db:       cfg.DB,
		config:   cfg,
		handlers: handlers,
	}
}

func (a *App) Start() error {
	// Middleware
	a.e.Use(middleware.Logger())
	a.e.Use(middleware.Recover())
	a.e.Use(middleware.CORS())

	// Setup routes
	a.setupRoutes()

	return a.e.Start(":" + a.config.DBConfig.Port)
}

func (a *App) setupRoutes() {
	// Initialize Kratos middleware
	kratosMiddleware := middleware.NewKratosMiddleware(a.handlers.Auth.KratosService)

	// Auth routes (public)
	auth := a.e.Group("/auth")
	auth.POST("/register", a.handlers.Auth.InitiateRegistration)
	auth.GET("/login", a.handlers.Auth.HandleLogin)
	auth.GET("/callback", a.handlers.Auth.HandleCallback)

	// Protected routes
	api := a.e.Group("/api")
	api.Use(kratosMiddleware.ValidateSession)

	// College-specific routes with role checks
	college := api.Group("/college/:collegeID")
	college.Use(kratosMiddleware.RequireCollege)

	// // Academic routes
	// academic := college.Group("/academic")
	// academic.GET("/attendance", a.handlers.Academic.ViewAttendance)
	// academic.POST("/attendance", a.handlers.Academic.MarkAttendance,
	// 	kratosMiddleware.RequireRole(RoleFaculty, RoleAdmin))

	// // Finance routes
	// finance := college.Group("/finance")
	// finance.GET("/fees", a.handlers.Finance.ViewFees)
	// finance.POST("/fees", a.handlers.Finance.ManageFees,
	// 	kratosMiddleware.RequireRole(RoleAdmin))

	// // Admin routes
	// admin := college.Group("/admin")
	// admin.Use(kratosMiddleware.RequireRole(RoleAdmin))
	// admin.GET("/reports", a.handlers.Admin.ViewReports)
	// admin.POST("/users", a.handlers.Admin.ManageUsers)
}
