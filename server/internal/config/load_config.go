package config

import (
	"log"

	"github.com/uptrace/bun"
)

type Config struct {
	DB         *bun.DB
	DBConfig   *DBConfig
	AuthConfig *AuthConfig
}

func NewConfig() (*Config, error) {
	dbConfig, err := LoadDatabaseConfig()
	if err != nil {
		return nil, err
	}
	db := LoadDatabase()
	authConfig, err := LoadAuthConfig()
	if err != nil {
		log.Fatal("error loading auth config")
	}
	
	cfg := &Config{
		DB:         db,
		DBConfig:   dbConfig,
		AuthConfig: authConfig,
	}

	return cfg, nil
}

func LoadConfig() (*Config, error) {
	return NewConfig()
}
