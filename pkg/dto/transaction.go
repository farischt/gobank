package dto

type CreateTransactionDTO struct {
	To     uint    `json:"to" binding:"required"`
	Amount float64 `json:"amount" binding:"required"`
}
