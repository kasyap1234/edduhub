package app

import (
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/v5"
	"github.com/uptrace/bun"
	"gorm.io/gorm"
	"eduhub/server/internal/config"
) 


type App struct {
	r *chi.Mux 
	db *bun.DB

	// zitadel auth 


}

func(a*App)New()*App{
	return &App{
		r : chi.NewRouter(),

		
	}

}

func (a*App)Start()(*config.Config,error){
	cfg,err:=config.LoadConfig()
	if err !=nil {
		return  nil,err
	}
	return cfg,nil 
}