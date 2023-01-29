package main

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/farischt/gobank/config"
	_ "github.com/lib/pq"
)

type Storage interface {
	GetAccount(int) (*Account, error)
	GetAllAccount() ([]*Account, error)
	CreateAccount(*CreateAccountDTO) error
	DeleteAccount(int) error
}

func GetPgConnectionStr() string {
	c := config.GetConfig()

	host := c.GetString("DB_HOST")
	user := c.GetString("DB_USER")
	password := c.GetString("DB_PASSWORD")
	name := c.GetString("DB_NAME")
	port := c.GetString("DB_PORT")

	return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, name, port)
}

type PgStorage struct {
	db *sql.DB
}

func NewPgStorage() (*PgStorage, error) {

	url := GetPgConnectionStr()
	db, err := sql.Open("postgres", url)

	if err != nil {
		return nil, err
	} else if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Printf("PostgreSQL Database up and running on port %d", config.GetConfig().GetInt("DB_PORT"))
	return &PgStorage{
		db: db,
	}, nil
}

// Service Layer
func (s *PgStorage) GetAccount(id int) (*Account, error) {
	query := `SELECT * FROM account WHERE id = $1`

	row := s.db.QueryRow(query, id)
	account, err := scanAccount(row)

	if err != nil {
		return nil, err
	}

	return account, nil
}

func (s *PgStorage) GetAllAccount() ([]*Account, error) {
	query := `SELECT * FROM account OFFSET $1 LIMIT $2`
	rows, err := s.db.Query(query, 0, 10)

	if err != nil {
		return nil, err
	}

	accounts, err := scanAccounts(rows)

	if err != nil {
		return nil, err
	}

	return accounts, nil
}

func (s *PgStorage) CreateAccount(account *CreateAccountDTO) error {
	query := `INSERT INTO account (first_name, last_name) VALUES ($1, $2)`
	_, err := s.db.Exec(
		query,
		account.FirstName,
		account.LastName,
	)
	return err
}

func (s *PgStorage) DeleteAccount(id int) error {
	query := `DELETE FROM account WHERE "id" = $1`
	_, err := s.db.Exec(query, id)
	return err
}

/*
Utils functions
*/
func scanAccount(row *sql.Row) (*Account, error) {
	a := new(Account)
	err := row.Scan(&a.ID, &a.FirstName, &a.LastName, &a.Balance, &a.CreatedAt, &a.UpdatedAt)
	return a, err
}

func scanAccounts(rows *sql.Rows) ([]*Account, error) {
	accounts := []*Account{}
	for rows.Next() {
		account := new(Account)

		err := rows.Scan(
			&account.ID,
			&account.FirstName,
			&account.LastName,
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
