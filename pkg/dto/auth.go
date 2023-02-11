package dto

type LoginDTO struct {
	AccountNumber uint   `json:"account_number" binding:"required"`
	Password      string `json:"password" binding:"required"`
}
