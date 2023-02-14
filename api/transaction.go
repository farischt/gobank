package api

import (
	"encoding/json"
	"net/http"

	"github.com/farischt/gobank/pkg/dto"
	"github.com/farischt/gobank/pkg/services"
)

type TransactionHandler struct {
	service *services.Service
}

func NewTransactionHandler(service *services.Service) *TransactionHandler {
	return &TransactionHandler{
		service: service,
	}
}

/*
HandleTransfer routes the request to the appropriate handler for /transfer endpoint.
*/
func (s *TransactionHandler) HandleTransfer(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		return s.createTransaction(w, r)
	default:
		return NewApiError(http.StatusMethodNotAllowed, "method_not_allowed")
	}
}

/* ------------------------------- Controller ------------------------------- */

/*
handleTransfer is the controller that handles the POST /transfer endpoint.
*/
func (s *TransactionHandler) createTransaction(w http.ResponseWriter, r *http.Request) error {
	data := new(dto.CreateTransactionDTO)

	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return NewApiError(http.StatusBadRequest, "invalid_request_body")
	}
	defer r.Body.Close()

	tokenId, err := GetTokenFromCookie(r)
	if err != nil {
		return err
	}

	// TODO: Useles check, since the token is already checked in the middleware
	token, err := s.service.Session.Get(tokenId)
	if err != nil {
		return NewApiError(http.StatusUnauthorized, "unauthorized")
	}

	err = s.service.Transaction.Transfer(token.AccountId, data)
	if err != nil {
		// TODO: Handle error to return appropriate error code
		return err
	}

	return WriteJSON(w, http.StatusCreated, NewApiResponse(http.StatusCreated, data, r))
}
