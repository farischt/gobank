package types

import (
	"time"

	"github.com/farischt/gobank/utils"
)

type Account struct {
	ID        uint      `db:"id"`
	UserID    uint      `db:"user_id"`
	Balance   []uint8   `db:"balance"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	User      *User
}

type SerializedAccount struct {
	ID        uint            `json:"id"`
	UserID    uint            `json:"user_id"`
	Balance   float64         `json:"balance,omitempty"`
	CreatedAt time.Time       `json:"created_at,omitempty"`
	UpdatedAt time.Time       `json:"updated_at,omitempty"`
	User      *SerializedUser `json:"user,omitempty"`
}

func (a *Account) Serialize() SerializedAccount {

	var serializedUser *SerializedUser

	if a.User != nil {
		s := a.User.Serialize()
		serializedUser = &s
	} else {
		serializedUser = nil
	}

	return SerializedAccount{
		ID:        a.ID,
		UserID:    a.UserID,
		Balance:   utils.Uint8ToFloat(a.Balance),
		CreatedAt: a.CreatedAt,
		UpdatedAt: a.UpdatedAt,
		User:      serializedUser,
	}
}

func SerializeAccounts(accounts []*Account) []SerializedAccount {
	var serializedAccounts []SerializedAccount
	for _, account := range accounts {
		serializedAccounts = append(serializedAccounts, account.Serialize())
	}
	return serializedAccounts
}
