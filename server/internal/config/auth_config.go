package config

import (
	"context"
	"fmt"
	"os"
	"net/url"
	"strings"
	"github.com/zitadel/zitadel-go/pkg/client/zitadel"
	"github.com/zitadel/zitadel-go/v3/pkg"
	"github.com/zitadel/oidc/v3/pkg/client"
	"github.com/zitadel/oidc/v3/pkg/oidc"
)

type AuthConfig struct {
	Domain      string
	Key         string
	ClientID    string
	ClientSecret string 
	RedirectURI string
	Port        string
	PostLogoutRedirectURI string 
	Client      *zitadel.Client
}

func LoadAuthConfig() (*AuthConfig, error) {
	config := &AuthConfig{
		Domain:      os.Getenv("ZITADEL_DOMAIN"),
		Key:         os.Getenv("ZITADEL_KEY"),
		ClientID:    os.Getenv("ZITADEL_CLIENT_ID"),
		ClientSecret: os.Getenv("ZITADEL_ClIENT_SECRET"),
		RedirectURI: os.Getenv("ZITADEL_REDIRECT_URI"),
		Port:        os.Getenv("PORT"),
	}

	ctx :=context.Background()
	// Initialize Zitadel client
	client, err :=client.New(
		ctx,
		&url.URL{Scheme: "https", Host: config.Domain}, // Production should be HTTPS
		client.WithClientID(config.ClientID),
		client.WithClientSecret(config.ClientSecret), // Consider fetching from env vars
		client.WithRedirectURL(config.RedirectURI),
		client.WithPostLogoutRedirectURI(config.PostLogoutRedirectURI),
		client.WithInsecureAllowHTTP(strings.Contains(config.Domain, "localhost")), // Dev only
		
	)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize Zitadel client: %w", err)
	}

	config.Client = client
	return config, nil
}
