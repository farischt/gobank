package services

import "github.com/farischt/gobank/pkg/store"

type Service struct {
	Account     AccountService
	User        UserService
	Transaction TransactionService
	Session     SessionService
}

func New(store store.Store) *Service {
	return &Service{
		Account:     NewAccountService(store),
		User:        NewUserService(store),
		Transaction: NewTransactionService(store),
		Session:     NewSessionService(store),
	}
}
