package store

import (
	"database/sql"
	"errors"

	"github.com/farischt/gobank/dto"
	"github.com/farischt/gobank/types"
)

type UserStore struct {
	db *sql.DB
}

func NewUser(db *sql.DB) *UserStore {
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

func (s *UserStore) GetUserByID(id uint) (*types.User, error) {
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

/*
scanUser is a helper function to scan a row into an User.
It takes a row and returns an Account and an error.
*/
func scanUser(row *sql.Row) (*types.User, error) {
	a := new(types.User)
	err := row.Scan(&a.ID, &a.FirstName, &a.LastName, &a.Email, &a.CreatedAt, &a.UpdatedAt)
	return a, err
}