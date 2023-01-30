package main

import (
	"encoding/json"
	"net/http"
)

/*
handleUser routes the request to the appropriate handler for /user endpoint.
*/
func (s *ApiServer) HandleUser(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		return s.handleCreateUser(w, r)
	default:
		return NewApiError(http.StatusMethodNotAllowed, "method_not_allowed")
	}
}

/* ------------------------------- Controller ------------------------------- */

/*
handleCreateUser is the controller that handles the POST /user endpoint.
*/
func (s *ApiServer) handleCreateUser(w http.ResponseWriter, r *http.Request) error {
	data := new(CreateUserDTO)

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

	exist, err := s.store.GetUserByEmail(data.Email)
	if err == nil && exist != nil {
		return NewApiError(http.StatusBadRequest, "email_already_exist")
	}

	err = s.store.CreateUser(data)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusCreated, NewApiResponse(http.StatusCreated, data, r))
}
