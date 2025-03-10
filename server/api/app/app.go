package app

import (
	"eduhub/server/api/handler"
	"eduhub/server/internal/config"
	"eduhub/server/internal/services/auth"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/uptrace/bun"
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
	authService := auth.NewAuthService(cfg.AuthConfig)

	// Initialize handlers
	handlers := handler.NewHandlers(authService)

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
	// Auth routes
	auth := a.e.Group("/auth")
	auth.POST("/register", a.handlers.Auth.RegisterUser)
	auth.GET("/login", a.handlers.Auth.LoginUser)
	// auth.GET("/callback", a.handlers.Auth.HandleCallback)

	// Protected routes
	api := a.e.Group("/api")
	api.Use(a.handlers.Auth.AuthService.VerifyTokenMiddleware)

	// Organization specific routes
	college := api.Group("/college/:orgID")
	college.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			orgID := c.Param("orgID")
			return a.handlers.Auth.AuthService.RequireOrganization(orgID)(next)(c)
		}
	})
}
