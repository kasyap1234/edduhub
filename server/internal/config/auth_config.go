package config

import (
	"context"
	"fmt"
	"os"

	"github.com/zitadel/zitadel-go/pkg/client/zitadel"
	"github.com/zitadel/zitadel-go/v3/pkg/zitadel"
)

type AuthConfig struct {
	Domain      string
	Key         string
	ClientID    string
	RedirectURI string
	Port        string
	Client      *zitadel.Client
}

func LoadAuthConfig() (*AuthConfig, error) {
	config := &AuthConfig{
		Domain:      os.Getenv("ZITADEL_DOMAIN"),
		Key:         os.Getenv("ZITADEL_KEY"),
		ClientID:    os.Getenv("ZITADEL_CLIENT_ID"),
		RedirectURI: os.Getenv("ZITADEL_REDIRECT_URI"),
		Port:        os.Getenv("PORT"),
	}

	// Initialize Zitadel client
	client, err := zitadel.New(
		context.Background(),
		zitadel.WithCustomDomain(config.Domain),
		zitadel.WithJWTProfileTokenSource(config.ClientID, []byte(config.Key)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Zitadel client: %w", err)
	}

	config.Client = client
	return config, nil
}
