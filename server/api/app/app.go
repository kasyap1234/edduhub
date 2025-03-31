package app

import (
	"eduhub/server/api/handler"
	"eduhub/server/internal/config"
	"eduhub/server/internal/middleware"
	"eduhub/server/internal/services"

	// "eduhub/server/internal/services/auth"

	// localmid "eduhub/server/internal/middleware"

	"github.com/labstack/echo/v4"
	echomid "github.com/labstack/echo/v4/middleware"
	"github.com/uptrace/bun"
)

type App struct {
	e          *echo.Echo
	db         *bun.DB
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
	handler.SetupRoutes(a.e, a.handlers, a.middleware)

	return a.e.Start(":" + a.config.DBConfig.Port)
}
