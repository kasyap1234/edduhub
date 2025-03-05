package config

import (
	"github.com/uptrace/bun"
)
type Config struct {
	DB *bun.DB 
	DBConfig DBConfig
	// Auth AuthConfig 

}

func NewConfig()(*Config,error){{
	dbConfig :=Start()
	db:=LoadDatabase()
	return &Config{
DB: db, 
DBConfig: dbConfig,
	},nil 
}}

func LoadConfig()(*Config,error){
return NewConfig()
}




