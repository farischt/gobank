package store

import (
	"database/sql"
	"errors"

	"github.com/farischt/gobank/dto"
	"github.com/farischt/gobank/types"
)

type AccountStore struct {
	db *sql.DB
}

func NewAccount(db *sql.DB) *AccountStore {
	return &AccountStore{db: db}
}

/*
GetAccount is a method to get an account by id.
It takes an id and returns an Account and an error.
*/
func (s *AccountStore) GetAccount(id uint) (*types.Account, error) {
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
func (s *AccountStore) GetAllAccount() ([]*types.Account, error) {
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
func (s *AccountStore) CreateAccount(account *dto.CreateAccountDTO) error {
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
func (s *AccountStore) DeleteAccount(id uint) error {
	query := `DELETE FROM account WHERE "id" = $1`
	_, err := s.db.Query(query, id)
	return err
}

/*
scanAccount is a helper function to scan a row into an Account.
It takes a row and returns an Account and an error.
*/
func scanAccount(row *sql.Row) (*types.Account, error) {
	a := new(types.Account)
	err := row.Scan(&a.ID, &a.Balance, &a.CreatedAt, &a.UpdatedAt, &a.UserID)
	return a, err
}

/*
scanAccounts is a helper function to scan rows into an array of Account.
It takes rows and returns an array of Account and an error.
*/
func scanAccounts(rows *sql.Rows) ([]*types.Account, error) {
	accounts := []*types.Account{}
	for rows.Next() {
		account := new(types.Account)

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