package types

import (
	"time"

	"github.com/farischt/gobank/utils"
)

type Transaction struct {
	ID        uint      `db:"id"`
	From      uint      `db:"from"`
	To        uint      `db:"to"`
	Amount    []uint8   `db:"amount"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type SerializedTransaction struct {
	ID        uint      `json:"id"`
	From      uint      `json:"from"`
	To        uint      `json:"to"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at" omitempty:"true"`
	UpdatedAt time.Time `json:"updated_at" omitempty:"true"`
}

func SerializeTransaction(t Transaction) SerializedTransaction {
	return SerializedTransaction{
		ID:        t.ID,
		From:      t.From,
		To:        t.To,
		Amount:    utils.Uint8ToFloat(t.Amount),
		CreatedAt: t.CreatedAt,
		UpdatedAt: t.UpdatedAt,
	}
}
