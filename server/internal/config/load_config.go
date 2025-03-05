package config

import (
	"github.com/uptrace/bun"
)
type Config struct {
	db *bun.DB 
	//auth 

}

func NewConfig()(*Config){{
	return &Config{

	}
}}
func LoadConfig(){


}