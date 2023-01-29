package main

/* ----------------------------- Authentication ----------------------------- */
type LoginDTO struct {
	AccountNumber uint `json:"account_number" binding:"required"`
}

/* ---------------------------------- User ---------------------------------- */

type CreateUserDTO struct {
	FirstName string `json:"first_name" binding:"required"`
	LastName  string `json:"last_name" binding:"required"`
	Email     string `json:"email" binding:"required"`
}

/* --------------------------------- Account -------------------------------- */

type CreateAccountDTO struct {
	UserID uint `json:"user_id" binding:"required"`
}

type UpdateAccountDTO CreateAccountDTO

/* --------------------------------- Transaction ---------------------------- */

type CreateTransactionDTO struct {
	To     uint    `json:"to" binding:"required"`
	Amount float64 `json:"amount" binding:"required"`
}
