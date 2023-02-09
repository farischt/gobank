package store

import (
	"fmt"

	"github.com/farischt/gobank/dto"
	"github.com/farischt/gobank/types"
	"github.com/jmoiron/sqlx"
)

type TransactionStore struct {
	db *sqlx.DB
}

func NewTransaction(db *sqlx.DB) *TransactionStore {
	return &TransactionStore{db: db}
}

func (s *TransactionStore) CreateTxn(from uint, data *dto.CreateTransactionDTO) error {
	query := `INSERT INTO transaction (from_id, to_id, amount) VALUES ($1, $2, $3)`
	_, err := s.db.Exec(
		query,
		from,
		data.To,
		data.Amount,
	)
	return err
}

func (s *TransactionStore) CreateTxnAndUpdateBalance(from *types.Account, to *types.Account, fromFinalBalance float64, toFinalBalance float64, data *dto.CreateTransactionDTO) error {
	// Start transaction
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}

	// defer rollback if error
	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	createTxnQuery := `INSERT INTO transaction (from_id, to_id, amount) VALUES ($1, $2, $3)`
	_, err = tx.Exec(
		createTxnQuery,
		from.ID,
		to.ID,
		data.Amount,
	)

	if err != nil {
		return fmt.Errorf("error creating transaction")
	}

	// Update balance
	updateFromBalanceQuery := `UPDATE account SET balance = $1 WHERE id = $2`
	_, err = tx.Exec(
		updateFromBalanceQuery,
		fromFinalBalance,
		from.ID,
	)

	if err != nil {
		return fmt.Errorf("error updating from account balance")
	}

	updateToBalanceQuery := `UPDATE account SET balance = $1 WHERE id = $2`
	_, err = tx.Exec(
		updateToBalanceQuery,
		toFinalBalance,
		to.ID,
	)

	if err != nil {
		return fmt.Errorf("error updating to account balance")
	}

	return nil
}
