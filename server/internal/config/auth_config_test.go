package config

import (
	"os"
	"reflect"
	"testing"
)

func TestLoadAuthConfig(t *testing.T) {
	// ✅ Set the required environment variables
	os.Setenv("KRATOS_PUBLIC_URL", "http://localhost:4433")
	os.Setenv("KRATOS_ADMIN_URL", "http://localhost:4434")
	os.Setenv("KRATOS_DOMAIN", "example.com")
	os.Setenv("PORT", "8080")

	// ✅ Ensure environment variables are cleared after the test
	defer func() {
		os.Unsetenv("KRATOS_PUBLIC_URL")
		os.Unsetenv("KRATOS_ADMIN_URL")
		os.Unsetenv("KRATOS_DOMAIN")
		os.Unsetenv("PORT")
	}()

	tests := []struct {
		name    string
		want    *AuthConfig
		wantErr bool
	}{
		{
			name: "Valid Auth Config",
			want: &AuthConfig{
				PublicURL: "http://localhost:4433",
				AdminURL:  "http://localhost:4434",
				Domain:    "example.com",
				Port:      "8080",
				College: CollegeConfig{
					RequireVerification: true,
					AllowedRoles:        []string{"admin", "faculty", "student"},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := LoadAuthConfig()

			if (err != nil) != tt.wantErr {
				t.Errorf("LoadAuthConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("LoadAuthConfig() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
