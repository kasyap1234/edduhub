package handler

import "eduhub/server/internal/services/auth"

type Handlers struct {
	Auth *AuthHandler
	// other handlers 
	// quiz handler 
	// fee handler 
	// attendance handler 
}

func NewHandlers(kratosService *auth.KratosService)*Handlers{
	return &Handlers{
		Auth: NewAuthHandler(kratosService),
		
	}
}
