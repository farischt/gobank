package config

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
)

const (
	PORT                      = "PORT"
	SESSION_COOKIE_NAME       = "SESSION_COOKIE_NAME"
	SESSION_COOKIE_EXPIRATION = "SESSION_COOKIE_EXPIRATION"
	HOST                      = "HOST"
	DB_HOST                   = "POSTGRES_HOSTNAME"
	DB_PORT                   = "POSTGRES_PORT"
	DB_USER                   = "POSTGRES_USER"
	DB_PASSWORD               = "POSTGRES_PASSWORD"
	DB_NAME                   = "POSTGRES_DB"
)

var config *viper.Viper
var dbConfig *viper.Viper

// Init is an exported method that takes the environment starts the viper
// (external lib) and returns the configuration struct.
func InitBaseConfig(env string) {
	var err error
	config = viper.New()

	config.AddConfigPath(".././")
	config.AddConfigPath("./")

	config.SetConfigType("env")
	config.SetConfigName(fmt.Sprintf(".env.%s", env))

	err = config.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
}

func InitDbConfig(env string) {
	var err error
	dbConfig = viper.New()

	dbConfig.AddConfigPath(".././")
	dbConfig.AddConfigPath("./")

	dbConfig.SetConfigType("env")
	dbConfig.SetConfigName(fmt.Sprintf(".env.%s.postgres", env))

	err = dbConfig.ReadInConfig()
	if err != nil {
		log.Fatal(err)
	}
}

// GetConfig is an exported method that returns the configuration struct.
func GetConfig() *viper.Viper {
	return config
}

func GetDbConfig() *viper.Viper {
	return dbConfig
}
