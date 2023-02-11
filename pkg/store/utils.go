package store

import (
	"fmt"

	"github.com/farischt/gobank/config"
)

/*
GetPgConnectionStr is a helper function to get the connection string to PostgreSQL.
*/
func getPgConnectionStr() string {
	c := config.GetDbConfig()

	host := c.GetString(config.DB_HOST)
	user := c.GetString(config.DB_USER)
	password := c.GetString(config.DB_PASSWORD)
	name := c.GetString(config.DB_NAME)
	port := c.GetInt(config.DB_PORT)

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable", host, user, password, name, port)
}
