package config

import (
	

	"github.com/uptrace/bun"
)

type Config struct {
	DB       *bun.DB
	DBConfig DBConfig
	AuthConfig *AuthConfig
}

func NewConfig() (*Config, error) {
	dbConfig := Start()
	db := LoadDatabase()
	authConfig:= AuthConfig()
	cfg := &Config{
		DB:       db,
		DBConfig: *dbConfig,
		AuthConfig: authConfig,
	}

	return cfg, nil
}

func LoadConfig() (*Config, error) {
	return NewConfig()
}
