package api

import (
	"encoding/json"
	"net/http"

	"github.com/farischt/gobank/dto"
	"github.com/farischt/gobank/store"
)

type UserHandler struct {
	store store.Store
}

func NewUserHandler(store store.Store) *UserHandler {
	return &UserHandler{store: store}
}

func (u *UserHandler) HandleUser(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		return u.handleCreateUser(w, r)
	default:
		return NewApiError(http.StatusMethodNotAllowed, "method_not_allowed")
	}
}

func (u *UserHandler) handleCreateUser(w http.ResponseWriter, r *http.Request) error {
	data := new(dto.CreateUserDTO)

	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return NewApiError(http.StatusBadRequest, "invalid_request_body")
	}
	defer r.Body.Close()

	if len(data.FirstName) == 0 {
		return NewApiError(http.StatusBadRequest, "empty_first_name")
	} else if len(data.LastName) == 0 {
		return NewApiError(http.StatusBadRequest, "empty_last_name")
	} else if len(data.Email) == 0 {
		return NewApiError(http.StatusBadRequest, "empty_email")
	}

	exist, err := u.store.User.GetUserByEmail(data.Email)
	if err == nil && exist != nil {
		return NewApiError(http.StatusBadRequest, "email_already_exist")
	}

	err = u.store.User.CreateUser(data)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusCreated, NewApiResponse(http.StatusCreated, data, r))
}
