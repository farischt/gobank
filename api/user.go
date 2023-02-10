package api

import (
	"encoding/json"
	"net/http"

	"github.com/farischt/gobank/dto"
	"github.com/farischt/gobank/services"
)

type UserHandler struct {
	service *services.Service
}

func NewUserHandler(service *services.Service) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

/*
HandleUser routes the request to the appropriate handler for /user endpoint.
*/
func (u *UserHandler) HandleUser(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		return u.createUser(w, r)
	default:
		return NewApiError(http.StatusMethodNotAllowed, "method_not_allowed")
	}
}

/*
HandleUniqueUser routes the request to the appropriate handler for /user/{id} endpoint.
*/
func (u *UserHandler) HandleUniqueUser(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return u.getUserById(w, r)
	default:
		return NewApiError(http.StatusMethodNotAllowed, "method_not_allowed")
	}
}

/*
createUser is the controller method that handles the POST /user endpoint.
*/
func (u *UserHandler) createUser(w http.ResponseWriter, r *http.Request) error {
	data := new(dto.CreateUserDTO)

	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return NewApiError(http.StatusBadRequest, "invalid_request_body")
	}
	defer r.Body.Close()

	err := u.service.User.Create(data)

	if err != nil {
		switch err.Error() {
		case "empty_first_name":
			return NewApiError(http.StatusBadRequest, err.Error())
		case "empty_last_name":
			return NewApiError(http.StatusBadRequest, err.Error())
		case "empty_email":
			return NewApiError(http.StatusBadRequest, err.Error())
		case "user_already_exists":
			return NewApiError(http.StatusBadRequest, err.Error())
		default:
			return err
		}
	}

	return WriteJSON(w, http.StatusCreated, NewApiResponse(http.StatusCreated, data, r))
}

/*
getUserById is the controller method that handles the GET /user/{id} endpoint.
*/
func (u *UserHandler) getUserById(w http.ResponseWriter, r *http.Request) error {
	id, err := GetIntParameter(r, "id")

	if err != nil {
		return NewApiError(http.StatusBadRequest, "missing_user_id")
	}

	user, err := u.service.User.Get(id)
	if err != nil {
		switch err.Error() {
		case "invalid_user_id":
			return NewApiError(http.StatusBadRequest, err.Error())
		default:
			return err
		}
	}

	return WriteJSON(w, http.StatusOK, NewApiResponse(http.StatusOK, user, r))
}
