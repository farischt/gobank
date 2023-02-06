package types

import (
	"time"
)

/* ---------------------------------- User ---------------------------------- */
type User struct {
	ID        uint      `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at" omitempty:"true"`
	UpdatedAt time.Time `json:"updated_at" omitempty:"true"`
}

/* --------------------------------- Account -------------------------------- */
type Account struct {
	ID        uint      `json:"id"`
	UserID    uint      `json:"user_id"`
	Balance   []uint8   `json:"balance" omitempty:"true"`
	CreatedAt time.Time `json:"created_at" omitempty:"true"`
	UpdatedAt time.Time `json:"updated_at" omitempty:"true"`
}

/* --------------------------------- Transaction ---------------------------- */
type Transaction struct {
	ID        uint      `json:"id"`
	From      uint      `json:"from"`
	To        uint      `json:"to"`
	Amount    float64   `json:"amount"`
	CreatedAt time.Time `json:"created_at" omitempty:"true"`
	UpdatedAt time.Time `json:"updated_at" omitempty:"true"`
}
