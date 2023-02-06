package dto

type LoginDTO struct {
	AccountNumber uint `json:"account_number" binding:"required"`
}
