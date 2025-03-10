package app

import (
	"eduhub/server/internal/config"

	"github.com/go-chi/chi/middleware"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/uptrace/bun"
)

type App struct {
    e *echo.Echo
	db     *bun.DB
	config *config.Config
	authConfig *config.AuthConfig
}


func (a *App) New() *App {
	authcfg,err :=config.LoadAuthConfig()
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	
	return &App{
		e : echo.New(),
		db:     cfg.DB,
		config: cfg,
		authConfig: authcfg,
	}

}

func(a*App)Start()error {
	a.e.Use(middleware.Logger())
	a.e.Use(middleware.Recoverer())
	a.e.Use(middleware.CORS())
	a.setupRoutes()
	return a.e.Start(":" + config.Port) 

}

func(a*App) setupRoutes(){
	a.e.GET("/auth/login",a.handleLogin)
	a.e.GET("/auth/register", a.RegisterUser)
}

