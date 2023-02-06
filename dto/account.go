package dto

type CreateAccountDTO struct {
	UserID uint `json:"user_id" binding:"required"`
}

type UpdateAccountDTO CreateAccountDTO
