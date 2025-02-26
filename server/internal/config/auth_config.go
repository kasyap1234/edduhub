package config

import ("github.com/go-chi/chi"
"github.com/go-chi/cors"
    "github.com/supertokens/supertokens-golang/supertokens"

)

type Router struct {
    Router *chi.Mux
}


func NewRouter()*Router{
	r :=chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
        AllowedOrigins: []string{"http://localhost:3000"},
        AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders: append([]string{"Content-Type"},
            supertokens.GetAllCORSHeaders()...),
        AllowCredentials: true,
    }))

    // Add SuperTokens middleware
    r.Use(supertokens.Middleware)

    return &Router{Router: r}
}
