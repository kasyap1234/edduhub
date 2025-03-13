package app

import (
	"eduhub/server/api/handler"
	"eduhub/server/internal/config"
	"eduhub/server/internal/services"
	"eduhub/server/internal/services/auth"

	"eduhub/server/internal/middleware"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/uptrace/bun"
)

type App struct {
	e        *echo.Echo
	db       *bun.DB
	config   *config.Config
	services *services.Services
	handlers *handler.Handlers
}

func New() *App {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}

	// Initialize auth service
	services :=services.NewServices(cfg)
	handlers:=handler.NewHandlers(services)



	// Initialize handlers
	

	return &App{
		e:        echo.New(),
		db:       cfg.DB,
		config:   cfg,
		services: services,
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

func (a*App)setupRoutes(){

	// auth routes  
	// protected college routes 
	// protected finance routes 
	

}