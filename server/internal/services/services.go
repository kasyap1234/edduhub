package services

import (
	"eduhub/server/internal/config"
	"eduhub/server/internal/services/auth"
)

type Services struct {
	Auth *auth.KratosService
	// Quiz *Quiz.QuizService 
	// Fee *Fee.FeeService 
}

func NewServices(cfg *config.Config)*Services{
	return &Services{
		Auth: auth.NewKratosService(),
		
	}

}