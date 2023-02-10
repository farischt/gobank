package store

import (
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

type Store struct {
	User         UserStorer
	Account      AccountStorer
	Transaction  TransactionStorer
	SessionToken SessionTokenStorer
}

func NewPostgres() (*Store, error) {
	url := getPgConnectionStr()
	db, err := sqlx.Connect("postgres", url)

	if err != nil {
		return nil, err
	} else if err := db.Ping(); err != nil {
		return nil, err
	}

	log.Println("Succesfully connected to postgres database")

	return &Store{
		User:         NewUser(db),
		Account:      NewAccount(db),
		Transaction:  NewTransaction(db),
		SessionToken: NewSessionToken(db),
	}, nil
}
