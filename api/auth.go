package api

import (
	"encoding/json"
	"net/http"

	"github.com/farischt/gobank/dto"
	"github.com/farischt/gobank/services"
)

type AuthenticationHandler struct {
	service *services.Service
}

func NewAuthenticationHandler(service *services.Service) *AuthenticationHandler {
	return &AuthenticationHandler{
		service: service,
	}
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
	defer r.Body.Close()

	token, err := h.service.Session.Create(data.AccountNumber)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, NewApiResponse(http.StatusOK, token, r))
}
