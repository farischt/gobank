package services

import (
	"fmt"

	"github.com/farischt/gobank/dto"
	"github.com/farischt/gobank/store"
	"github.com/farischt/gobank/types"
)

type TransactionService interface {
	Transfer(senderId uint, data *dto.CreateTransactionDTO) error
}

type transactionService struct {
	store store.Store
}

func NewTransactionService(store store.Store) TransactionService {
	return &transactionService{
		store: store,
	}
}

func (t *transactionService) Create() error {
	return nil
}

func (t *transactionService) Transfer(senderId uint, data *dto.CreateTransactionDTO) error {

	if data.Amount <= 0 {
		return fmt.Errorf("invalid_amount")
	} else if data.To <= 0 {
		return fmt.Errorf("invalid_to_account_id")
	} else if data.To == senderId {
		return fmt.Errorf("cannot_transfer_to_yourself")
	}

	sender, err := t.store.Account.GetAccount(senderId)
	if err != nil {
		return err
	}

	s := sender.Serialize()
	if !t.HasEnoughBalance(&s, data.Amount) {
		return fmt.Errorf("insufficient_balance")
	}

	recipient, err := t.store.Account.GetAccount(data.To)
	if err != nil {
		return err
	}

	r := recipient.Serialize()
	s.Balance -= data.Amount
	r.Balance += data.Amount

	return t.store.Transaction.CreateTxnAndUpdateBalance(sender, recipient, s.Balance, r.Balance, data)

}

func (t *transactionService) HasEnoughBalance(account *types.SerializedAccount, amount float64) bool {
	return account.Balance >= amount
}
