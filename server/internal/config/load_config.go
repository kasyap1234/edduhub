package config

import (
	"os"

	"github.com/uptrace/bun"
)

type Config struct {
	DB       *bun.DB
	DBConfig DBConfig
	Auth     struct {
		Domain      string
		Key         string
		ClientID    string
		RedirectURI string
		Port        string
	}
}

func NewConfig() (*Config, error) {
	dbConfig := Start()
	db := LoadDatabase()

	cfg := &Config{
		DB:       db,
		DBConfig: *dbConfig,
		Auth: struct {
			Domain      string
			Key         string
			ClientID    string
			RedirectURI string
			Port        string
		}{
			Domain:      os.Getenv("domain"),
			Key:         os.Getenv("key"),
			ClientID:    os.Getenv("clientid"),
			RedirectURI: os.Getenv("redirecturi"),
			Port:        os.Getenv("port"),
		},
	}

	return cfg, nil
}

func LoadConfig() (*Config, error) {
	return NewConfig()
}
