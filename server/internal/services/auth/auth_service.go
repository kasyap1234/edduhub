package auth

import (
	"context"
	"eduhub/server/internal/config"
	"fmt"
	"strings"

	"github.com/zitadel/zitadel-go/v3/pkg/client/management"
	"github.com/zitadel/zitadel-go/v3/pkg/client/zitadel"
)

type AuthService struct {
	client     *zitadel.Client
	authConfig *config.AuthConfig
}

type RegisterUserParams struct {
	Email     string
	Password  string
	FirstName string
	LastName  string
	OrgID     string
}

func NewAuthService(cfg *config.AuthConfig) *AuthService {
	return &AuthService{
		client:     cfg.Client,
		authConfig: cfg,
	}
}

func (s *AuthService) GetLoginURL() string {
	return fmt.Sprintf("%s/oauth/v2/authorize?"+
		"client_id=%s&"+
		"redirect_uri=%s&"+
		"response_type=code&"+
		"scope=%s",
		s.authConfig.Domain,
		s.authConfig.ClientID,
		s.authConfig.RedirectURI,
		strings.Join(s.authConfig.Scopes, " "),
	)
}

func (s *AuthService) RegisterUser(ctx context.Context, params RegisterUserParams) error {
	// Check if organization exists
	org, err := s.client.ManagementAPI().GetOrganizationByID(ctx, params.OrgID)
	if err != nil {
		return fmt.Errorf("invalid organization: %w", err)
	}

	// Create user in the specific organization
	_, err = s.client.ManagementAPI().CreateUser(ctx, &management.CreateUserRequest{
		Username: params.Email,
		Password: params.Password,
		Profile: &management.Profile{
			FirstName: params.FirstName,
			LastName:  params.LastName,
			Email:     params.Email,
		},
		OrganizationID: org.ID,
	})

	return err
}

func (s *AuthService) ExchangeCodeForToken(ctx context.Context, code string) (*zitadel.Token, error) {
	return s.client.ExchangeAuthCode(ctx, code, s.authConfig.RedirectURI)
}

func (s *AuthService) VerifyToken(ctx context.Context, token string) (*zitadel.TokenClaims, error) {
	return s.client.VerifyAccessToken(ctx, token)
}

func (s *AuthService) ValidateOrganization(ctx context.Context, orgID string) error {
	_, err := s.client.ManagementAPI().GetOrganizationByID(ctx, orgID)
	return err
}
