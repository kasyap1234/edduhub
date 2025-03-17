package app

import (
	"eduhub/server/api/handler"
	"eduhub/server/internal/config"
	"eduhub/server/internal/services"
	// "eduhub/server/internal/services/auth"

	// localmid "eduhub/server/internal/middleware"

	"github.com/labstack/echo/v4"
	echomid "github.com/labstack/echo/v4/middleware"
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
	a.e.Use(echomid.Logger())
	a.e.Use(echomid.Recover())
	a.e.Use(echomid.CORS())

	// Setup routes
	a.setupRoutes()

	return a.e.Start(":" + a.config.DBConfig.Port)
}

func (a*App)setupRoutes(){

	// auth routes  

	// protected college routes 
	// protected finance routes 
auth :=a.e.Group("/auth")

auth.POST("/register",a.handlers.Auth.InitiateRegistration)
auth.POST("/auth/register/complete",a.handlers.Auth.HandleRegistration)
auth.GET("/login",a.handlers.Auth.HandleLogin)
auth.GET("/callback",a.handlers.Auth.HandleCallback)


}