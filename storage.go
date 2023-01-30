package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"github.com/farischt/gobank/config"
	_ "github.com/lib/pq"
)

type Storage interface {
	// User
	CreateUser(*CreateUserDTO) error
	GetUserByEmail(string) (*User, error)
	GetUserBydID(uint) (*User, error)

	// Account
	GetAccount(uint) (*Account, error)
	GetAllAccount() ([]*Account, error)
	CreateAccount(*CreateAccountDTO) error
	DeleteAccount(uint) error

	// Transaction
	CreateTxn(uint, *CreateTransactionDTO) error
	CreateTxnAndUpdateBalance(from *Account, to *Account, fromFinalBalance float64, toFinalBalance float64, data *CreateTransactionDTO) error
}

type PgStorage struct {
	db *sql.DB
}

func NewPgStorage() (*PgStorage, error) {
	url := getPgConnectionStr()
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

/* ---------------------------------- User ---------------------------------- */

/*
CreateUser is a method to create a user.
It takes a CreateUserDTO and returns an error.
*/
func (s *PgStorage) CreateUser(input *CreateUserDTO) error {
	query := `INSERT INTO "user" (first_name, last_name, email) VALUES ($1, $2, $3)`
	_, err := s.db.Exec(
		query,
		input.FirstName,
		input.LastName,
		input.Email,
	)

	return err
}

/*
GetUserByEmail is a method to get a user by email.
It takes an email and returns a User and an error.
*/
func (s *PgStorage) GetUserByEmail(email string) (*User, error) {
	query := `SELECT * FROM "user" WHERE email = $1`

	row := s.db.QueryRow(query, email)
	user, err := scanUser(row)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user_not_found")
		}

		return nil, err
	}

	return user, nil
}

func (s *PgStorage) GetUserBydID(id uint) (*User, error) {
	query := `SELECT * FROM "user" WHERE id = $1`

	row := s.db.QueryRow(query, id)
	user, err := scanUser(row)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user_not_found")
		}

		return nil, err
	}

	return user, nil
}

/* --------------------------------- Account -------------------------------- */

/*
GetAccount is a method to get an account by id.
It takes an id and returns an Account and an error.
*/
func (s *PgStorage) GetAccount(id uint) (*Account, error) {
	// query := `SELECT a.*, u.first_name, u.last_name, u.email FROM account AS a LEFT JOIN "user" AS u ON a.user_id = u.id WHERE a.id = $1`
	query := `SELECT * FROM account WHERE id = $1`

	row := s.db.QueryRow(query, id)
	account, err := scanAccount(row)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("account_not_found")
		}

		return nil, err
	}

	return account, nil
}

/*
GetAllAccount is a method to get all accounts.
It returns an array of Account and an error.
*/
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

/*
CreateAccount is a method to create an account.
It takes a CreateAccountDTO and returns an error.
*/
func (s *PgStorage) CreateAccount(account *CreateAccountDTO) error {
	query := `INSERT INTO account (user_id) VALUES ($1)`
	_, err := s.db.Exec(
		query,
		account.UserID,
	)
	return err
}

/*
DeleteAccount is a method to delete an account by id.
It takes an id and returns an error.
*/
func (s *PgStorage) DeleteAccount(id uint) error {
	query := `DELETE FROM account WHERE "id" = $1`
	_, err := s.db.Query(query, id)
	return err
}

/* --------------------------------- Transaction ---------------------------- */

func (s *PgStorage) CreateTxn(from uint, data *CreateTransactionDTO) error {
	query := `INSERT INTO transaction (from_id, to_id, amount) VALUES ($1, $2, $3)`
	_, err := s.db.Exec(
		query,
		from,
		data.To,
		data.Amount,
	)
	return err
}

func (s *PgStorage) CreateTxnAndUpdateBalance(from *Account, to *Account, fromFinalBalance float64, toFinalBalance float64, data *CreateTransactionDTO) error {
	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// defer rollback if error
	defer func() {
		if err != nil {
			tx.Rollback()
		} else {
			tx.Commit()
		}
	}()

	createTxnQuery := `INSERT INTO transaction (from_id, to_id, amount) VALUES ($1, $2, $3)`
	_, err = tx.Exec(
		createTxnQuery,
		from.ID,
		to.ID,
		data.Amount,
	)

	if err != nil {
		return fmt.Errorf("error creating transaction")
	}

	// Update balance
	updateFromBalanceQuery := `UPDATE account SET balance = $1 WHERE id = $2`
	_, err = tx.Exec(
		updateFromBalanceQuery,
		fromFinalBalance,
		from.ID,
	)

	if err != nil {
		return fmt.Errorf("error updating from account balance")
	}

	updateToBalanceQuery := `UPDATE account SET balance = $1 WHERE id = $2`
	_, err = tx.Exec(
		updateToBalanceQuery,
		toFinalBalance,
		to.ID,
	)

	if err != nil {
		return fmt.Errorf("error updating to account balance")
	}

	return nil
}
