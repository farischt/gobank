package services

import (
	"fmt"

	"github.com/farischt/gobank/dto"
	"github.com/farischt/gobank/store"
	"github.com/farischt/gobank/types"
)

type AccountService interface {
	Get(id uint, withUser bool) (*types.SerializedAccount, error)
	GetAll() ([]*types.SerializedAccount, error)
	Create(data *dto.CreateAccountDTO) error
}

type accountService struct {
	store store.Store
}

func NewAccountService(store store.Store) AccountService {
	return &accountService{
		store: store,
	}
}

func (a *accountService) Get(id uint, withUser bool) (*types.SerializedAccount, error) {
	var acc *types.Account
	var err error

	if id <= 0 {
		return nil, fmt.Errorf("invalid_account_id")
	}

	if withUser {
		acc, err = a.store.Account.GetAccountWithUser(id)
	} else {
		acc, err = a.store.Account.GetAccount(id)
	}

	if err != nil {
		return nil, err
	}

	s := acc.Serialize()

	return &s, nil
}

func (a *accountService) GetAll() ([]*types.SerializedAccount, error) {
	accounts, err := a.store.Account.GetAllAccount()
	if err != nil {
		return nil, err
	}

	var serializedAccounts []*types.SerializedAccount

	for _, acc := range accounts {
		s := acc.Serialize()
		serializedAccounts = append(serializedAccounts, &s)
	}

	return serializedAccounts, nil
}

func (a *accountService) Create(data *dto.CreateAccountDTO) error {

	if data.UserID <= 0 {
		return fmt.Errorf("invalid_user_id")
	}

	// Check if user exists
	_, err := a.store.User.GetUserByID(data.UserID)
	if err != nil {
		return err
	}

	return a.store.Account.CreateAccount(data)
}
