package store

import (
	"github.com/farischt/gobank/types"

	"github.com/farischt/gobank/dto"
)

type UserStorer interface {
	CreateUser(*dto.CreateUserDTO) error
	GetUserByEmail(string) (*types.User, error)
	GetUserByID(id uint) (*types.User, error)
}


type AccountStorer interface {
	GetAccount(id uint) (*types.Account, error)
	GetAllAccount() ([]*types.Account, error)
	CreateAccount(account *dto.CreateAccountDTO) error
	DeleteAccount(id uint) error
}


type TransactionStorer interface {
	CreateTxn(from uint, data *dto.CreateTransactionDTO) error
    CreateTxnAndUpdateBalance(from *types.Account, to *types.Account, fromFinalBalance float64, toFinalBalance float64, data *dto.CreateTransactionDTO) error
}


// type Storer interface {
// 	user UserStorer
// 	account AccountStorer
// 	transaction TransactionStorer
// }

