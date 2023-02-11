package dto

type CreateAccountDTO struct {
	UserID   uint   `json:"user_id" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UpdateAccountDTO CreateAccountDTO
