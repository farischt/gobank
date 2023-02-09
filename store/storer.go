package store

import (
	"github.com/farischt/gobank/types"

	"github.com/farischt/gobank/dto"
)

type UserStorer interface {
	CreateUser(input *dto.CreateUserDTO) error
	GetUserByEmail(email string) (*types.User, error)
	GetUserByID(id uint) (*types.User, error)
}

type AccountStorer interface {
	GetAccount(id uint) (*types.Account, error)
	GetAllAccount() ([]*types.Account, error)
	GetAccountWithUser(id uint) (*types.Account, error)
	CreateAccount(account *dto.CreateAccountDTO) error
	DeleteAccount(id uint) error
}

type TransactionStorer interface {
	CreateTxn(from uint, data *dto.CreateTransactionDTO) error
	CreateTxnAndUpdateBalance(from *types.Account, to *types.Account, fromFinalBalance float64, toFinalBalance float64, data *dto.CreateTransactionDTO) error
}

type SessionTokenStorer interface {
	CreateSessionToken(accountId uint) (string, error)
	GetSessionToken(token string) (*types.SessionToken, error)
	DeleteSessionToken(token string) error
	IsValidSessionToken(token string) (uint, bool)
}
