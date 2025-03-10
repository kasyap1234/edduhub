package config

import (
	"context"
	"fmt"
	"os"

	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel"
)

type AuthConfig struct {
	Domain         string
	ClientID       string
	ClientSecret   string
	RedirectURI    string
	Key            string
	Port           string
	OrganizationID string
	Scopes         []string
	// Client         *zitadel.Client
	Client *zitadel.Client
}


func LoadAuthConfig() (*AuthConfig, error) {
	config := &AuthConfig{
		Domain:         os.Getenv("ZITADEL_DOMAIN"),
		ClientID:       os.Getenv("ZITADEL_CLIENT_ID"),
		ClientSecret:   os.Getenv("ZITADEL_CLIENT_SECRET"),
		RedirectURI:    os.Getenv("ZITADEL_REDIRECT_URI"),
		Key:            os.Getenv("ZITADEL_KEY"),
		Port:           os.Getenv("PORT"),
		OrganizationID: os.Getenv("ZITADEL_ORG_ID"),
		Scopes:         []string{"openid", "profile", "email"},
	}

	// Validate required fields
	if config.Domain == "" || config.ClientID == "" || config.ClientSecret == "" {
		return nil, fmt.Errorf("missing required Zitadel configuration")
	}

	// Initialize Zitadel client
	client, err := zitadel.New(
		context.Background(),
		zitadel.WithCustomDomain(config.Domain),
		zitadel.WithClientID(config.ClientID),
		zitadel.WithClientSecret(config.ClientSecret),
		zitadel.WithJWTProfileTokenSource(config.ClientID, []byte(config.Key)),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Zitadel client: %w", err)
	}

	config.Client = client
	return config, nil
}
