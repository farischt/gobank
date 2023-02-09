package api

import (
	"encoding/json"
	"net/http"

	"github.com/farischt/gobank/dto"
	"github.com/farischt/gobank/store"
)

type AuthenticationHandler struct {
	store store.Store
}

func NewAuthenticationHandler(store store.Store) *AuthenticationHandler {
	return &AuthenticationHandler{store: store}
}

/*
HandleLogin routes the request to the appropriate handler for /auth/login endpoint.
*/
func (h *AuthenticationHandler) HandleLogin(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		return h.createToken(w, r)
	default:
		return NewApiError(http.StatusMethodNotAllowed, "method_not_allowed")
	}
}

/*
handleCreateToken is the controller that handles the POST /auth/login endpoint.
It creates a new token for the user.
*/
func (h *AuthenticationHandler) createToken(w http.ResponseWriter, r *http.Request) error {
	data := new(dto.LoginDTO)

	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return NewApiError(http.StatusBadRequest, "invalid_request_body")
	}
	//defer r.Body.Close()

	if data.AccountNumber <= 0 {
		return NewApiError(http.StatusBadRequest, "missing_account_number")
	}

	// Check if the account exists
	a, err := h.store.Account.GetAccount(data.AccountNumber)
	if err != nil {
		if err.Error() == "account_not_found" {
			return NewApiError(http.StatusNotFound, err.Error())
		}
		return err
	}

	// TODO Check if the password is correct
	token, err := h.store.SessionToken.CreateSessionToken(a.ID)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, NewApiResponse(http.StatusOK, token, r))
}
