package config

import (
	"os"
	"testing"
)

func TestLoadDatabaseConfig(t *testing.T) {
	tests := []struct {
		name        string
		envVars     map[string]string
		expectError bool
		expected    *DBConfig
	}{
		{
			name: "valid configuration",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_PORT":     "5432",
				"DB_USER":     "testuser",
				"DB_PASSWORD": "testpass",
				"DB_NAME":     "testdb",
				"DB_SSLMODE":  "require",
			},
			expectError: false,
			expected: &DBConfig{
				Host:     "localhost",
				Port:     "5432",
				User:     "testuser",
				Password: "testpass",
				DBName:   "testdb",
				SSLMode:  "require",
			},
		},
		{
			name: "missing required fields",
			envVars: map[string]string{
				"DB_HOST": "localhost",
				"DB_PORT": "5432",
			},
			expectError: true,
		},
		{
			name: "invalid port number",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_PORT":     "invalid",
				"DB_USER":     "testuser",
				"DB_PASSWORD": "testpass",
				"DB_NAME":     "testdb",
			},
			expectError: true,
		},
		{
			name: "default ssl mode",
			envVars: map[string]string{
				"DB_HOST":     "localhost",
				"DB_PORT":     "5432",
				"DB_USER":     "testuser",
				"DB_PASSWORD": "testpass",
				"DB_NAME":     "testdb",
			},
			expectError: false,
			expected: &DBConfig{
				Host:     "localhost",
				Port:     "5432",
				User:     "testuser",
				Password: "testpass",
				DBName:   "testdb",
				SSLMode:  "disable",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			os.Clearenv()
			for k, v := range tt.envVars {
				os.Setenv(k, v)
			}

			config, err := LoadDatabaseConfig()

			if tt.expectError && err == nil {
				t.Error("expected error but got none")
			}

			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}

			if !tt.expectError && tt.expected != nil {
				if config.Host != tt.expected.Host {
					t.Errorf("expected Host %s, got %s", tt.expected.Host, config.Host)
				}
				if config.Port != tt.expected.Port {
					t.Errorf("expected Port %s, got %s", tt.expected.Port, config.Port)
				}
				if config.User != tt.expected.User {
					t.Errorf("expected User %s, got %s", tt.expected.User, config.User)
				}
				if config.Password != tt.expected.Password {
					t.Errorf("expected Password %s, got %s", tt.expected.Password, config.Password)
				}
				if config.DBName != tt.expected.DBName {
					t.Errorf("expected DBName %s, got %s", tt.expected.DBName, config.DBName)
				}
				if config.SSLMode != tt.expected.SSLMode {
					t.Errorf("expected SSLMode %s, got %s", tt.expected.SSLMode, config.SSLMode)
				}
			}
		})
	}
}

func TestBuildDSN(t *testing.T) {
	tests := []struct {
		name     string
		config   DBConfig
		expected string
	}{
		{
			name: "valid config",
			config: DBConfig{
				Host:     "localhost",
				Port:     "5432",
				User:     "testuser",
				Password: "testpass",
				DBName:   "testdb",
				SSLMode:  "disable",
			},
			expected: "postgres://testuser:testpass@localhost:5432/testdb?sslmode=disable",
		},
		{
			name: "with special characters",
			config: DBConfig{
				Host:     "test.host",
				Port:     "5432",
				User:     "test@user",
				Password: "test:pass",
				DBName:   "test-db",
				SSLMode:  "require",
			},
			expected: "postgres://test@user:test:pass@test.host:5432/test-db?sslmode=require",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := buildDSN(tt.config)
			if result != tt.expected {
				t.Errorf("expected DSN %s, got %s", tt.expected, result)
			}
		})
	}
}
