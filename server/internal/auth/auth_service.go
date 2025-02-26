package auth

import (
	"github.com/go-chi/chi"
	"github.com/kasyap1234/eduhub-backend/server/internal/config"
	"github.com/kasyap1234/eduhub-backend/server/internal/repository"
	"github.com/supertokens/supertokens-golang/recipe/emailpassword"
)

type SuperTokensAuth struct {
	config *config.SuperTokensConfig
	userRepo repository.UserRepository
}

func NewAuthService(config *config.SuperTokensConfig,userRepo repository.UserRepository)*SuperTokensAuth{
	return &SuperTokensAuth{
		config: config, 
		userRepo : userRepo,
	}
}

