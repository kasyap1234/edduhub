package handler

import "eduhub/server/internal/services/auth"

type Handlers struct {
	Auth *AuthHandler
	// other handlers 
}

func NewHandlers(authService *auth.AuthService)*Handlers{
	return &Handlers{
		Auth: NewAuthHandler(authService),
	}
}
