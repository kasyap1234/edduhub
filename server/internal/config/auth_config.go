package config

import (
	"fmt"
	"os"
)

type AuthConfig struct {
	PublicURL string
	AdminURL  string
	Domain    string
	Port      string
	College   CollegeConfig
}

type CollegeConfig struct {
	RequireVerification bool
	AllowedRoles        []string
}

func LoadAuthConfig() (*AuthConfig, error) {
	config := &AuthConfig{
		PublicURL: os.Getenv("KRATOS_PUBLIC_URL"),
		AdminURL:  os.Getenv("KRATOS_ADMIN_URL"),
		Domain:    os.Getenv("KRATOS_DOMAIN"),
		Port:      os.Getenv("PORT"),
		College: CollegeConfig{
			RequireVerification: true,
			AllowedRoles:        []string{"admin", "faculty", "student"},
		},
	}

	// Validate required fields
	if config.PublicURL == "" || config.AdminURL == "" {
		return nil, fmt.Errorf("missing required Kratos configuration")
	}

	return config, nil
}
