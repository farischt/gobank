package store

import (
	"database/sql"

	_ "github.com/lib/pq"
)


type Store struct {
	User UserStore 
	Account AccountStore
	Transaction TransactionStore
}

func NewPgStore() (*Store, error ){
	url := getPgConnectionStr()
	db, err := sql.Open("postgres", url)

	if err != nil {
		return nil, err
	} else if err := db.Ping(); err != nil {
		return nil, err
	}

	return &Store{
		User: *NewUser(db),
		Account: *NewAccount(db),
		Transaction: *NewTransaction(db),
	}, nil
}