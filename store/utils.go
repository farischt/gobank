package store

import (
	"fmt"

	"github.com/farischt/gobank/config"
)

/*
GetPgConnectionStr is a helper function to get the connection string to PostgreSQL.
*/
func getPgConnectionStr() string {
	c := config.GetConfig()

	host := c.GetString("DB_HOST")
	user := c.GetString("DB_USER")
	password := c.GetString("DB_PASSWORD")
	name := c.GetString("DB_NAME")
	port := c.GetString("DB_PORT")

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, name, port)
}

