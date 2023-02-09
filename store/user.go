package store

import (
	"database/sql"
	"errors"

	"github.com/farischt/gobank/dto"
	"github.com/farischt/gobank/types"
	"github.com/jmoiron/sqlx"
)

type UserStore struct {
	db *sqlx.DB
}

func NewUser(db *sqlx.DB) *UserStore {
	return &UserStore{db: db}
}

/*
CreateUser is a method to create a user.
It takes a CreateUserDTO and returns an error.
*/
func (s *UserStore) CreateUser(input *dto.CreateUserDTO) error {
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
func (s *UserStore) GetUserByEmail(email string) (*types.User, error) {
	query := `SELECT * FROM "user" WHERE email = $1`

	user := new(types.User)
	err := s.db.QueryRowx(query, email).StructScan(user)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("user_not_found")
		}

		return nil, err
	}

	return user, nil
}

/*
GetUserByID is a method to get a user by id.
It takes an id and returns a User and an error.
*/
func (s *UserStore) GetUserByID(id uint) (*types.User, error) {
	query := `SELECT * FROM "user" WHERE id = $1`

	user := new(types.User)
	err := s.db.QueryRowx(query, id).StructScan(user)

	if err != nil {

		if err == sql.ErrNoRows {
			return nil, errors.New("user_not_found")
		}

		return nil, err
	}

	return user, nil
}
