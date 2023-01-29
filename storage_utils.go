package main

import (
	"database/sql"
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

/*
scanAccount is a helper function to scan a row into an Account.
It takes a row and returns an Account and an error.
*/
func scanAccount(row *sql.Row) (*Account, error) {
	a := new(Account)
	err := row.Scan(&a.ID, &a.Balance, &a.CreatedAt, &a.UpdatedAt, &a.UserID)
	return a, err
}

/*
scanAccounts is a helper function to scan rows into an array of Account.
It takes rows and returns an array of Account and an error.
*/
func scanAccounts(rows *sql.Rows) ([]*Account, error) {
	accounts := []*Account{}
	for rows.Next() {
		account := new(Account)

		err := rows.Scan(
			&account.ID,
			&account.UserID,
			&account.Balance,
			&account.CreatedAt,
			&account.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
	}

	return accounts, nil
}

/*
scanUser is a helper function to scan a row into an User.
It takes a row and returns an Account and an error.
*/
func scanUser(row *sql.Row) (*User, error) {
	a := new(User)
	err := row.Scan(&a.ID, &a.FirstName, &a.LastName, &a.Email, &a.CreatedAt, &a.UpdatedAt)
	return a, err
}
