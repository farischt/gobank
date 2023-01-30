package main

import (
	"encoding/json"
	"net/http"
)

/*
handleAccount routes the request to the appropriate handler for /account endpoint.
*/
func (s *ApiServer) HandleAccount(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetAccounts(w, r)
	case "POST":
		return s.handleCreateAccount(w, r)
	default:
		return NewApiError(http.StatusMethodNotAllowed, "method_not_allowed")
	}
}

/*
handleUniqueAccount routes the request to the appropriate handler for /account/{id} endpoint.
*/
func (s *ApiServer) HandleUniqueAccount(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetAccount(w, r)
	case "DELETE":
		return s.handleDeleteAccount(w, r)
	default:
		return NewApiError(http.StatusMethodNotAllowed, "method_not_allowed")
	}
}

/* ------------------------------- Controller ------------------------------- */

/*
handleGetAccounts is the controller that handles the GET /account endpoint.
*/
func (s *ApiServer) handleGetAccounts(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAllAccount()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, NewApiResponse(http.StatusOK, accounts, r))
}

/*
handleCreateAccount is the controller that handles the POST /account endpoint.
*/
func (s *ApiServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	data := new(CreateAccountDTO)

	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return NewApiError(http.StatusBadRequest, "invalid_request_body")
	}
	defer r.Body.Close()

	if data.UserID <= 0 {
		return NewApiError(http.StatusBadRequest, "invalid_user_id")
	}

	// check if user exists
	user, err := s.store.GetUserBydID(data.UserID)
	if err != nil {
		if err.Error() == "user_not_found" {
			return NewApiError(http.StatusNotFound, err.Error())
		}
		return err
	}
	data.UserID = user.ID

	err = s.store.CreateAccount(data)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusCreated, NewApiResponse(http.StatusCreated, data, r))
}

/*
handleGetAccount is the controller that handles the GET /account/{id} endpoint.
*/
func (s *ApiServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := GetIntParameter(r, "id")
	if err != nil {
		return err
	}

	a, err := s.store.GetAccount(id)
	if err != nil {
		if err.Error() == "account_not_found" {
			return NewApiError(http.StatusNotFound, err.Error())
		}
		return err
	}

	return WriteJSON(w, http.StatusOK, NewApiResponse(http.StatusOK, a, r))
}

/*
handleDeleteAccount is the controller that handles the DELETE /account/{id} endpoint.
*/
func (s *ApiServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	id, err := GetIntParameter(r, "id")
	if err != nil {
		return err
	}

	account, err := s.store.GetAccount(id)
	if err != nil {
		if err.Error() == "account_not_found" {
			return NewApiError(http.StatusNotFound, err.Error())
		}
		return err
	}

	err = s.store.DeleteAccount(id)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, NewApiResponse(http.StatusOK, map[string]uint{"deleted": account.ID}, r))
}
