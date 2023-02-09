package store

import (
	"database/sql"
	"errors"
	"time"

	"github.com/farischt/gobank/dto"
	"github.com/farischt/gobank/types"
	"github.com/jmoiron/sqlx"
)

type AccountStore struct {
	db *sqlx.DB
}

func NewAccount(db *sqlx.DB) *AccountStore {
	return &AccountStore{db: db}
}

/*
GetAccount is a method to get an account by id.
It takes an id and returns an Account and an error.
*/
func (s *AccountStore) GetAccount(id uint) (*types.Account, error) {
	// query := `SELECT a.*, u.first_name, u.last_name, u.email FROM account AS a LEFT JOIN "user" AS u ON a.user_id = u.id WHERE a.id = $1`
	query := `SELECT * FROM account WHERE id = $1`

	account := new(types.Account)

	err := s.db.QueryRowx(query, id).StructScan(account)
	//account, err := scanAccount(row)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("account_not_found")
		}

		return nil, err
	}

	return account, nil
}

/*
GetAccountWithUser is a method to get an account by id with the corresponding user.
It takes an id and returns an Account and an error.
*/
func (s *AccountStore) GetAccountWithUser(id uint) (*types.Account, error) {

	query := `SELECT a.*, a.balance, u.id AS uid , u.first_name, u.last_name, u.email, u.created_at AS ucreated_at, u.updated_at AS uupadted_at FROM account AS a LEFT JOIN "user" AS u ON a.user_id = u.id WHERE a.id = $1`

	var result struct {
		ID        uint      `db:"id"`
		UserId    uint      `db:"user_id"`
		Balance   []uint8   `db:"balance"`
		CreatedAt time.Time `db:"created_at"`
		UpdatedAt time.Time `db:"updated_at"`
		// User relation
		UID        uint      `db:"uid"`
		FirstName  string    `db:"first_name"`
		LastName   string    `db:"last_name"`
		Email      string    `db:"email"`
		UCreatedAt time.Time `db:"ucreated_at"`
		UUpdatedAt time.Time `db:"uupadted_at"`
	}

	err := s.db.Get(&result, query, id)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("account_not_found")
		}

		return nil, err
	}

	account := &types.Account{
		ID:        result.ID,
		UserID:    result.UserId,
		Balance:   result.Balance,
		CreatedAt: result.CreatedAt,
		UpdatedAt: result.UpdatedAt,
		User: &types.User{
			ID:        result.UID,
			FirstName: result.FirstName,
			LastName:  result.LastName,
			Email:     result.Email,
			CreatedAt: result.UCreatedAt,
			UpdatedAt: result.UUpdatedAt,
		},
	}

	return account, nil
}

/*
GetAllAccount is a method to get all accounts.
It returns an array of Account and an error.
*/
func (s *AccountStore) GetAllAccount() ([]*types.Account, error) {
	query := `SELECT * FROM account OFFSET $1 LIMIT $2`
	rows, err := s.db.Queryx(query, 0, 10)

	if err != nil {
		return nil, err
	}

	accounts := []*types.Account{}

	for rows.Next() {
		account := new(types.Account)
		err := rows.StructScan(account)
		if err != nil {
			return nil, err
		}

		accounts = append(accounts, account)
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
