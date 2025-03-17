package config

import (
	"github.com/uptrace/bun"
	"log"
)

type Config struct {
	DB         *bun.DB
	DBConfig   DBConfig
	AuthConfig *AuthConfig
}

func NewConfig() (*Config, error) {
	dbConfig := LoadDatabaseConfig()
	db := LoadDatabase()
	authConfig, err := LoadAuthConfig()
	if err != nil {
		log.Fatal("error loading auth config")
	}

	cfg := &Config{
		DB:         db,
		DBConfig:   *dbConfig,
		AuthConfig: authConfig,
	}

	return cfg, nil
}

func LoadConfig() (*Config, error) {
	return NewConfig()
}
