package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

const (
	PORT         = "PORT"
	TOKEN_SECRET = "TOKEN_SECRET"
	TOKEN_NAME   = "TOKEN_NAME"
	HOST         = "HOST"
	DB_HOST      = "DB_HOST"
	DB_PORT      = "DB_PORT"
	DB_USER      = "DB_USER"
	DB_PASSWORD  = "DB_PASSWORD"
	DB_NAME      = "DB_NAME"
)

var config *viper.Viper

// Init is an exported method that takes the environment starts the viper
// (external lib) and returns the configuration struct.
func Init(env string) {
	var err error
	config = viper.New()
	config.SetConfigType("json")
	config.SetConfigName(fmt.Sprintf("config.%s", env))
	config.AddConfigPath("../config/")
	config.AddConfigPath("config/")
	err = config.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
}

// GetConfig is an exported method that returns the configuration struct.
func GetConfig() *viper.Viper {
	return config
}
