package handler

import (
	"os"

	"github.com/kasyap1234/eduhub-backend/server/internal/auth"
	"github.com/kasyap1234/eduhub-backend/server/internal/config"
	"github.com/kasyap1234/eduhub-backend/server/internal/repository"
	"gorm.io/gorm"
)

func authHandler() {
	supertokensConfig := &config.SuperTokensConfig{
		ConnectionURI: os.Getenv("SUPERTOKENS_URI"),
		APIKey: os.Getenv("SUPERTOKENS_API_KEY"),
		AppName: "EDUHUB",
		APIDomain : "http://localhost:8080",
		WebDomain: "http://localhost:3000",
	}

	config.InitSuperTokens(*supertokensConfig)

	userRepo := repository.NewUserRepository(db)
	authService :=auth.NewAuthService(supertokensConfig,userRepo)


}