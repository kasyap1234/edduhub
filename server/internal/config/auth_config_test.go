package config

import (
	"os"
	"testing"
)

func TestLoadAuthConfig(t *testing.T) {
	tests := []struct {
		name          string
		envVars       map[string]string
		expectError   bool
		expectedValue *AuthConfig
	}{
		{
			name: "valid configuration",
			envVars: map[string]string{
				"KRATOS_PUBLIC_URL": "http://public.example.com",
				"KRATOS_ADMIN_URL":  "http://admin.example.com",
				"KRATOS_DOMAIN":     "example.com",
				"PORT":              "8080",
			},
			expectError: false,
			expectedValue: &AuthConfig{
				PublicURL: "http://public.example.com",
				AdminURL:  "http://admin.example.com",
				Domain:    "example.com",
				Port:      "8080",
				College: CollegeConfig{
					RequireVerification: true,
					AllowedRoles:        []string{"admin", "faculty", "student"},
				},
			},
		},
		{
			name: "missing public URL",
			envVars: map[string]string{
				"KRATOS_ADMIN_URL": "http://admin.example.com",
				"KRATOS_DOMAIN":    "example.com",
				"PORT":             "8080",
			},
			expectError: true,
		},
		{
			name: "missing admin URL",
			envVars: map[string]string{
				"KRATOS_PUBLIC_URL": "http://public.example.com",
				"KRATOS_DOMAIN":     "example.com",
				"PORT":              "8080",
			},
			expectError: true,
		},
		{
			name: "optional fields can be empty",
			envVars: map[string]string{
				"KRATOS_PUBLIC_URL": "http://public.example.com",
				"KRATOS_ADMIN_URL":  "http://admin.example.com",
			},
			expectError: false,
			expectedValue: &AuthConfig{
				PublicURL: "http://public.example.com",
				AdminURL:  "http://admin.example.com",
				College: CollegeConfig{
					RequireVerification: true,
					AllowedRoles:        []string{"admin", "faculty", "student"},
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clear environment before each test
			os.Clearenv()

			// Set environment variables for test
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			config, err := LoadAuthConfig()

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !tt.expectError && tt.expectedValue != nil {
				if config.PublicURL != tt.expectedValue.PublicURL {
					t.Errorf("expected PublicURL %s, got %s", tt.expectedValue.PublicURL, config.PublicURL)
				}
				if config.AdminURL != tt.expectedValue.AdminURL {
					t.Errorf("expected AdminURL %s, got %s", tt.expectedValue.AdminURL, config.AdminURL)
				}
				if config.Domain != tt.expectedValue.Domain {
					t.Errorf("expected Domain %s, got %s", tt.expectedValue.Domain, config.Domain)
				}
				if config.Port != tt.expectedValue.Port {
					t.Errorf("expected Port %s, got %s", tt.expectedValue.Port, config.Port)
				}
				if config.College.RequireVerification != tt.expectedValue.College.RequireVerification {
					t.Errorf("expected RequireVerification %v, got %v", tt.expectedValue.College.RequireVerification, config.College.RequireVerification)
				}
				if len(config.College.AllowedRoles) != len(tt.expectedValue.College.AllowedRoles) {
					t.Errorf("expected AllowedRoles length %d, got %d", len(tt.expectedValue.College.AllowedRoles), len(config.College.AllowedRoles))
				}
			}
		})
	}
}
