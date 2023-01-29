package main

import (
	"time"
)

type CreateAccountDTO struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
}

type UpdateAccountDTO CreateAccountDTO

type Account struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`

	Balance   int64     `json:"balance" omitempty:"true"`
	CreatedAt time.Time `json:"created_at" omitempty:"true"`
	UpdatedAt time.Time `json:"updated_at" omitempty:"true"`
}

type SecuredAccount struct {
	ID        int    `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
