package app

import (
	"eduhub/server/api/handler"
	"eduhub/server/internal/config"
	"eduhub/server/internal/middleware"
	"eduhub/server/internal/repository"
	"eduhub/server/internal/services"

	"github.com/labstack/echo/v4"
	echomid "github.com/labstack/echo/v4/middleware"
)

type App struct {
	e          *echo.Echo
	db         *repository.DB
	config     *config.Config
	services   *services.Services
	handlers   *handler.Handlers
	middleware *middleware.Middleware
}

func New() *App {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	// Initialize auth service
	services := services.NewServices(cfg)
	handlers := handler.NewHandlers(services)
	// repos := repository.NewRepository(cfg.DB)
	mid := middleware.NewMiddleware(services)

	return &App{
		e:          echo.New(),
		db:         cfg.DB,
		config:     cfg,
		services:   services,
		handlers:   handlers,
		middleware: mid,
	}
}

func (a *App) Start() error {
	// Middleware

	a.e.Use(echomid.Logger())
	a.e.Use(echomid.Recover())
	a.e.Use(echomid.CORS())

	// Setup routes
	handler.SetupRoutes(a.e, a.handlers, a.middleware.Auth)

	return a.e.Start(":" + a.config.AppPort)
}
