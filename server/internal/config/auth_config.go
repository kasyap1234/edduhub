package config

import "os"

type AuthConfig struct {
	Domain      string
	Key         string
	ClientID    string
	RedirectURI string
	Port        string
}

func LoadAuthConfig() *AuthConfig {
	config := &AuthConfig{
		Domain: os.Getenv("domain"),
		Key: os.Getenv("key"),
		ClientID: os.Getenv("client"),
		RedirectURI: os.Getenv("redirecturi"),
		Port: os.Getenv("port"),

	}
	return config; 
}