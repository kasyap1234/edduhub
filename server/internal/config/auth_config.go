package config

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/zitadel/oidc/v3/pkg/client"
)

type AuthConfig struct {
	Domain                string
	Key                   string
	ClientID              string
	ClientSecret          string
	RedirectURI           string
	Port                  string
	PostLogoutRedirectURI string
	Client                *client.Client
	Scopes                []string
}

func LoadAuthConfig() (*AuthConfig, error) {
	domain := os.Getenv("ZITADEL_DOMAIN")
	if domain == "" {
		return nil, fmt.Errorf("ZITADEL_DOMAIN is required")
	}

	config := &AuthConfig{
		Domain:       domain,
		Key:          os.Getenv("ZITADEL_KEY"),
		ClientID:     os.Getenv("ZITADEL_CLIENT_ID"),
		ClientSecret: os.Getenv("ZITADEL_CLIENT_SECRET"),
		RedirectURI:  os.Getenv("ZITADEL_REDIRECT_URI"),
		Port:         os.Getenv("PORT"),
		Scopes:       []string{"openid", "profile", "email"},
	}

	// Validate required fields
	if config.ClientID == "" || config.ClientSecret == "" || config.RedirectURI == "" {
		return nil, fmt.Errorf("missing required configuration")
	}

	ctx := context.Background()
	// Initialize OIDC client
	oidcClient, err := client.New(
		ctx,
		&url.URL{Scheme: "https", Host: config.Domain},
		client.WithClientID(config.ClientID),
		client.WithClientSecret(config.ClientSecret),
		client.WithRedirectURL(config.RedirectURI),
		client.WithScopes(config.Scopes...),
		client.WithPostLogoutRedirectURI(config.PostLogoutRedirectURI),
		client.WithInsecureAllowHTTP(strings.Contains(config.Domain, "localhost")),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize OIDC client: %w", err)
	}

	config.Client = oidcClient
	return config, nil
}
