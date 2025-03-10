package config

import (
	

	"github.com/uptrace/bun"
)

type Config struct {
	DB       *bun.DB
	DBConfig DBConfig
	// AuthConfig *AuthConfig
}

func NewConfig() (*Config, error) {
	dbConfig := Start()
	db := LoadDatabase()

	cfg := &Config{
		DB:       db,
		DBConfig: *dbConfig,
		
	}

	return cfg, nil
}

func LoadConfig() (*Config, error) {
	return NewConfig()
}
