package services

import (
	"fmt"

	"github.com/farischt/gobank/pkg/dto"
	"github.com/farischt/gobank/pkg/store"
	"github.com/farischt/gobank/pkg/types"
)

type UserService interface {
	Create(data *dto.CreateUserDTO) error
	Get(id uint) (*types.SerializedUser, error)
}

type userService struct {
	store store.Store
}

func NewUserService(store store.Store) UserService {
	return &userService{
		store: store,
	}
}

func (u *userService) Get(id uint) (*types.SerializedUser, error) {
	if id <= 0 {
		return nil, fmt.Errorf("invalid_user_id")
	}

	user, err := u.store.User.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	s := user.Serialize()
	return &s, nil
}

func (u *userService) Create(data *dto.CreateUserDTO) error {
	if len(data.FirstName) == 0 {
		return fmt.Errorf("empty_first_name")
	} else if len(data.LastName) == 0 {
		return fmt.Errorf("empty_last_name")
	} else if len(data.Email) == 0 {
		return fmt.Errorf("empty_email")
	}

	exist, err := u.store.User.GetUserByEmail(data.Email)
	if err == nil && exist != nil {
		return fmt.Errorf("user_already_exist")
	}

	err = u.store.User.CreateUser(data)
	return err
}
