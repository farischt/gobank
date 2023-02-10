package types

import "time"

type SessionToken struct {
	ID        string    `db:"id"`
	AccountId uint      `db:"account_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type SerializedSessionToken struct {
	ID        string    `json:"id"`
	AccountId uint      `json:"account_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (s *SessionToken) Serialize() *SerializedSessionToken {
	return &SerializedSessionToken{
		ID:        s.ID,
		AccountId: s.AccountId,
		CreatedAt: s.CreatedAt,
		UpdatedAt: s.UpdatedAt,
	}
}
