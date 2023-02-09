package api

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/farischt/gobank/dto"
	"github.com/farischt/gobank/store"
)

type TransactionHandler struct {
	store store.Store
}

func NewTransactionHandler(store store.Store) *TransactionHandler {
	return &TransactionHandler{store: store}
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

	if data.Amount <= 0 {
		return NewApiError(http.StatusBadRequest, "invalid_amount")
	} else if data.To == 0 {
		return NewApiError(http.StatusBadRequest, "invalid_to_account_id")
	}

	id := GetAuthenticatedAccountId(r, s.store)

	fromAccount, err := s.store.Account.GetAccount(*id)
	if err != nil {
		if err.Error() == "account_not_found" {
			return NewApiError(http.StatusNotFound, "from_account_not_found")
		}
		return err
	}

	// Check if the to account exists
	toAccount, err := s.store.Account.GetAccount(data.To)
	if err != nil {
		if err.Error() == "account_not_found" {
			return NewApiError(http.StatusNotFound, "to_account_not_found")
		}
		return err
	}

	if fromAccount.ID == toAccount.ID {
		return NewApiError(http.StatusBadRequest, "cannot_transfer_to_same_account")
	}

	balance, _ := strconv.ParseFloat(string(fromAccount.Balance), 64)
	if balance < data.Amount {
		return NewApiError(http.StatusBadRequest, "insufficient_balance")
	}

	// Update the balance of the from account
	fromBalance := balance - data.Amount

	toAccountBalance, _ := strconv.ParseFloat(string(toAccount.Balance), 64)
	toBalance := toAccountBalance + data.Amount

	err = s.store.Transaction.CreateTxnAndUpdateBalance(fromAccount, toAccount, fromBalance, toBalance, data)
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusCreated, NewApiResponse(http.StatusCreated, data, r))
}
