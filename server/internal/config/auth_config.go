package config

import (
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/supertokens/supertokens-golang/supertokens"
)

// RouterConfig holds all router-specific configuration
type RouterConfig struct {
	AllowedOrigins     []string
	AllowedMethods     []string
	BaseAllowedHeaders []string
	AllowCredentials   bool
}

// Router wraps chi router with additional configuration
type Router struct {
	*chi.Mux
	config RouterConfig
}

// NewDefaultConfig returns default router configuration
func NewDefaultConfig() RouterConfig {
	return RouterConfig{
		AllowedOrigins:     []string{"http://localhost:3000"},
		AllowedMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		BaseAllowedHeaders: []string{"Content-Type"},
		AllowCredentials:   true,
	}
}

// NewRouter creates a new router with the given configuration
func NewRouter(config RouterConfig) *Router {
	r := &Router{
		Mux:    chi.NewRouter(),
		config: config,
	}
	r.setupMiddleware()
	return r
}

// setupMiddleware configures all necessary middleware
func (r *Router) setupMiddleware() {
	r.Use(r.corsMiddleware())
	r.Use(supertokens.Middleware)
}

// corsMiddleware returns configured CORS middleware
func (r *Router) corsMiddleware() func(next http.Handler) http.Handler {
	return cors.Handler(cors.Options{
		AllowedOrigins:   r.config.AllowedOrigins,
		AllowedMethods:   r.config.AllowedMethods,
		AllowedHeaders:   append(r.config.BaseAllowedHeaders, supertokens.GetAllCORSHeaders()...),
		AllowCredentials: r.config.AllowCredentials,
	})
}
