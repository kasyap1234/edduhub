package config

import (
	"os"

	"eduhub/server/internal/repository"
)

type Config struct {
	DB         *repository.DB
	DBConfig   *DBConfig
	AuthConfig *AuthConfig
	AppPort    string
}

func NewConfig() (*Config, error) {
	dbConfig, err := LoadDatabaseConfig()
	if err != nil {
		return nil, err
	}
	db := LoadDatabase()
	authConfig, err := LoadAuthConfig()
	if err != nil {
		return nil, err
	}

	AppPort := os.Getenv("APP_PORT")
	cfg := &Config{
		DB:         db,
		DBConfig:   dbConfig,
		AuthConfig: authConfig,
		AppPort:    AppPort,
	}

	return cfg, nil
}

func LoadConfig() (*Config, error) {
	return NewConfig()
}
