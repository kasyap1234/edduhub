package auth

import (
	"context"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel"
	"eduhub/server/internal/config"
)

type AuthService struct {
	client *zitadel.Client
	authConfig *config.AuthConfig

}


func NewAuthService(cfg *config.AuthConfig)*AuthService{
	return &AuthService{
		client : cfg.Client, 
		authConfig: cfg,
	}
}

func(s *AuthService)GetLoginURL()string {
	
}