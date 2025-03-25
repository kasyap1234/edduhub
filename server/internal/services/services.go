package services

import (
	"eduhub/server/internal/config"
	"eduhub/server/internal/services/auth"
)

type Services struct {
 Auth  auth.AuthService
	// Quiz *Quiz.QuizService 
	// Fee *Fee.FeeService 
}

func NewServices(cfg *config.Config)*Services{
	kratosService :=auth.NewKratosService()
	ketoService :=auth.NewKetoService()
	authService :=auth.NewAuthService(kratosService,ketoService)
	return &Services{
		Auth: authService, 
		
	}

}