package api

import (
	"encoding/json"
	"net/http"

	"github.com/farischt/gobank/dto"
	"github.com/farischt/gobank/services"
)

type AccountHandler struct {
	service *services.Service
}

func NewAccountHandler(service *services.Service) *AccountHandler {
	return &AccountHandler{
		service: service,
	}
}

/*
HandleAccount routes the request to the appropriate handler for /account endpoint.
*/
func (s *AccountHandler) HandleAccount(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.getAccounts(w, r)
	case "POST":
		return s.createAccount(w, r)
	default:
		return NewApiError(http.StatusMethodNotAllowed, "method_not_allowed")
	}
}

/*
handleUniqueAccount routes the request to the appropriate handler for /account/{id} endpoint.
*/
func (s *AccountHandler) HandleUniqueAccount(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.getAccount(w, r)
	default:
		return NewApiError(http.StatusMethodNotAllowed, "method_not_allowed")
	}
}

/* ------------------------------- Controller ------------------------------- */

/*
getAccounts is the controller method that handles the GET /account endpoint.
*/
func (s *AccountHandler) getAccounts(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.service.Account.GetAll()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, NewApiResponse(http.StatusOK, accounts, r))
}

/*
createAccount is the controller method that handles the POST /account endpoint.
*/
func (s *AccountHandler) createAccount(w http.ResponseWriter, r *http.Request) error {
	data := new(dto.CreateAccountDTO)

	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return NewApiError(http.StatusBadRequest, "invalid_request_body")
	}
	defer r.Body.Close()

	err := s.service.Account.Create(data)
	if err != nil {
		if err.Error() == "user_not_found" {
			return NewApiError(http.StatusBadRequest, "user_not_found")
		}
		return err
	}

	return WriteJSON(w, http.StatusCreated, NewApiResponse(http.StatusCreated, data, r))
}

/*
getAccount is the controller method that handles the GET /account/{id} endpoint.
*/
func (s *AccountHandler) getAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := GetIntParameter(r, "id")
	if err != nil {
		return NewApiError(http.StatusBadRequest, "missing_account_id")
	}

	param := r.URL.Query()
	_, exist := param["user"]

	a, err := s.service.Account.Get(id, exist)

	if err != nil {
		if err.Error() == "account_not_found" {
			return NewApiError(http.StatusNotFound, err.Error())
		} else if err.Error() == "invalid_account_owner" {
			return NewApiError(http.StatusBadRequest, err.Error())
		}
		return err
	}

	return WriteJSON(w, http.StatusOK, NewApiResponse(http.StatusOK, a, r))
}
