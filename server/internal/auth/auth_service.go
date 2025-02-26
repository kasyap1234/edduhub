package auth

import "github.com/go-chi/chi"

type AuthHandler struct {
	router chi.Router
	jwtSecret []byte 
	userStore UserStore
}