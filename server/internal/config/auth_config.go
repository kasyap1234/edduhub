// package config

// import (
// 	"context"
// 	"fmt"
// 	"os"

// 	"github.com/aws/aws-sdk-go/aws/client"
// 	"github.com/coreos/go-oidc/v3/oidc"
// 	"github.com/zitadel/zitadel-go/v3/pkg/client"
// 	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel"
// 	"golang.org/x/oauth2"
// )

// type AuthConfig struct {
// 	Domain         string
// 	ClientID       string
// 	ClientSecret   string
// 	RedirectURI    string
// 	Provider       *oidc.Provider
// 	Key            string
// 	Port           string
// 	OrganizationID string
// 	IssuerURL      string
// 	Scopes         []string
// 	Client         *zitadel.Client       // Zitadel client
// 	Verifier       *oidc.IDTokenVerifier // OIDC verifier
// 	OAuth2Config   *oauth2.Config        // OAuth2 config
// }

// func LoadAuthConfig() (*AuthConfig, error) {
// 	config := &AuthConfig{
// 		Domain:         os.Getenv("ZITADEL_DOMAIN"),
// 		ClientID:       os.Getenv("ZITADEL_CLIENT_ID"),
// 		ClientSecret:   os.Getenv("ZITADEL_CLIENT_SECRET"),
// 		RedirectURI:    os.Getenv("ZITADEL_REDIRECT_URI"),
// 		Key:            os.Getenv("ZITADEL_KEY"),
// 		Port:           os.Getenv("PORT"),
// 		OrganizationID: os.Getenv("ZITADEL_ORG_ID"),
// 		IssuerURL:      os.Getenv("ISSUER_URL"),
// 		Scopes:         []string{"openid", "profile", "email", "role"},
// 	}

// 	// Validate required fields
// 	if config.Domain == "" || config.ClientID == "" || config.ClientSecret == "" {
// 		return nil, fmt.Errorf("missing required Zitadel configuration")
// 	}

// 	ctx := context.Background()

// 	// Initialize OIDC provider
// 	provider, err := oidc.NewProvider(ctx, config.Domain)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to initialize OIDC provider: %w", err)
// 	}
// 	config.Provider = provider

// 	// Configure OAuth2
// 	config.OAuth2Config = &oauth2.Config{
// 		ClientID:     config.ClientID,
// 		ClientSecret: config.ClientSecret,
// 		RedirectURL:  config.RedirectURI,
// 		Endpoint:     provider.Endpoint(),
// 		Scopes:       config.Scopes,
// 	}

// 	// Configure ID Token verifier
// 	config.Verifier = provider.Verifier(&oidc.Config{
// 		ClientID: config.ClientID,
// 	})

// 	// Initialize Zitadel client
// 	client, err := client.New(
// 		ctx,
// 		zitadel.New(*&config.Domain),
// 		client.WithAuth(client.DefaultServiceUserAuthentication(*&config.Key,oidc.ScopeOpenID,client.ScopeZitadelAPI))
		
// 	)
// 	if err != nil {
// 		return nil, fmt.Errorf("failed to initialize Zitadel client: %w", err)
// 	}
// 	config.Client = client

// 	return config, nil
// }
package config

