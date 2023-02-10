package store

import (
	"database/sql"
	"errors"
	"time"

	"github.com/farischt/gobank/types"
	"github.com/jmoiron/sqlx"
)

type SessionTokenStore struct {
	db *sqlx.DB
}

func NewSessionToken(db *sqlx.DB) *SessionTokenStore {
	return &SessionTokenStore{
		db: db,
	}
}

/*
CreateSessionToken creates a new session token for the given account id.
It returns the token id and an error if any.
*/
func (s *SessionTokenStore) CreateSessionToken(accountId uint) (*types.SessionToken, error) {
	token := new(types.SessionToken)
	query := `INSERT INTO session_token (account_id) VALUES ($1) RETURNING *`
	// _, _ = s.db.NamedQuery(query, accountId)

	err := s.db.QueryRowx(query, accountId).StructScan(token)

	if err != nil {
		return nil, err
	}

	return token, nil
}

/*
GetSessionToken returns the session token for the given token id.
It returns an error if the token is not found.
*/
func (s *SessionTokenStore) GetSessionToken(token string) (*types.SessionToken, error) {

	query := `SELECT * FROM session_token WHERE id = $1`

	st := new(types.SessionToken)
	err := s.db.Get(st, query, token)

	if err != nil {
		if err == sql.ErrNoRows {
			return st, errors.New("session_token_not_found")
		}

		return st, err
	}

	return st, nil
}

/*
DeleteSessionToken deletes the session token for the given token id.
It returns an error if the token is not found.
*/
func (s *SessionTokenStore) DeleteSessionToken(token string) error {
	query := `DELETE FROM session_token WHERE id = $1`
	_, err := s.db.Exec(query, token)
	return err
}

// TODO: This method should available at the api level
func (s *SessionTokenStore) IsValidSessionToken(token string) (uint, bool) {
	st, err := s.GetSessionToken(token)
	if err != nil {
		return 0, false
	}

	elapsed := time.Since(st.CreatedAt)
	return st.AccountId, elapsed <= time.Second*1000
}
