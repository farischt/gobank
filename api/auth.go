package api

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/farischt/gobank/config"
	"github.com/farischt/gobank/pkg/dto"
	"github.com/farischt/gobank/pkg/services"
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
		return h.login(w, r)
	default:
		return NewApiError(http.StatusMethodNotAllowed, "method_not_allowed")
	}
}

func (h *AuthenticationHandler) HandleLogout(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		return h.logout(w, r)
	default:
		return NewApiError(http.StatusMethodNotAllowed, "method_not_allowed")
	}
}

/*
login is the controller that handles the POST /auth/login endpoint.
It creates a new token for the user.
*/
func (h *AuthenticationHandler) login(w http.ResponseWriter, r *http.Request) error {
	data := new(dto.LoginDTO)

	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return NewApiError(http.StatusBadRequest, "invalid_request_body")
	}
	defer r.Body.Close()

	token, err := h.service.Session.Create(data.AccountNumber, data.Password)

	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     config.GetConfig().GetString(config.SESSION_COOKIE_NAME),
		Value:    token.ID,
		Path:     "/",
		HttpOnly: true,
		MaxAge:   config.GetConfig().GetInt(config.SESSION_COOKIE_EXPIRATION), // Seconds
	})

	return WriteJSON(w, http.StatusOK, NewApiResponse(http.StatusOK, token, r))
}

func (h *AuthenticationHandler) logout(w http.ResponseWriter, r *http.Request) error {
	tokenId, err := GetTokenFromCookie(r)
	if err != nil {
		return err
	}

	err = h.service.Session.Delete(tokenId)

	if err != nil {
		return err
	}

	http.SetCookie(w, &http.Cookie{
		Name:     config.GetConfig().GetString(config.SESSION_COOKIE_NAME),
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Expires:  time.Now(),
	})

	return WriteJSON(w, http.StatusOK, NewApiResponse(http.StatusOK, nil, r))
}
