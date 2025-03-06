package config

import (
	"os"

	"github.com/uptrace/bun"
)
type Config struct {
	DB *bun.DB 
	DBConfig DBConfig
	Auth AuthConfig


}

func NewConfig()(*Config,error){{
	dbConfig :=Start()
	db:=LoadDatabase()
	// authConfig:=AuthConfig{
	// 	Domain: os.Getenv("domain"),
	// 	Key: os.Getenv("key"),
	// 	ClientID: os.Getenv("clientid"),
	// 	RedirectURI: os.Getenv("redirecturi"),
	// 	Port: os.Getenv("port"),

	// }
	authConfig:=LoadAuthConfig()
	return &Config{
DB: db, 
DBConfig: *dbConfig,
Auth : *authConfig,
	},nil 
}}

func LoadConfig()(*Config,error){
return NewConfig()
}




