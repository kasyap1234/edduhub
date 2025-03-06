package app

import (
	"eduhub/server/internal/config"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5"
	"github.com/uptrace/bun"
	
)

type App struct {
	r      *chi.Mux
	db     *bun.DB
	config *config.Config
}

func (a *App) New() *App {
	cfg, err := config.LoadConfig()
	if err != nil {
		panic(err)
	}
	return &App{
		r:      chi.NewRouter(),
		db:     cfg.DB,
		config: cfg,
	}

}

func (a *App) Start() (*config.Config, error) {
	cfg, err := config.LoadConfig()

	if err != nil {
		return nil, err
	}
	return cfg, nil
}
