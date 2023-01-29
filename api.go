package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

/*
ApiServer is the API server.
*/
type ApiServer struct {
	listenAddr string
	store      Storage
}

/*
NewApiServer creates a new instance of API server.
*/
func NewApiServer(l string, s Storage) *ApiServer {
	return &ApiServer{
		listenAddr: l,
		store:      s,
	}
}

/*
Start starts the API server.
*/
func (s *ApiServer) Start() {
	router := mux.NewRouter()

	router.HandleFunc("/user", makeHTTPFunc(s.handleUser))
	router.HandleFunc("/auth/login", withoutAuth(makeHTTPFunc(s.handleLogin)))
	router.HandleFunc("/account", makeHTTPFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", makeHTTPFunc(s.handleUniqueAccount))
	router.HandleFunc("/transfer", withAuth(makeHTTPFunc(s.handleTransfer)))

	log.Println("Server up and running on port ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

/* ------------------------------- Handlers ------------------------------ */

/*
handleUser routes the request to the appropriate handler for /user endpoint.
*/
func (s *ApiServer) handleUser(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		return s.handleCreateUser(w, r)
	default:
		return NewApiError(http.StatusMethodNotAllowed, "method_not_allowed")
	}
}

func (s *ApiServer) handleLogin(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		return s.handleCreateToken(w, r)
	default:
		return NewApiError(http.StatusMethodNotAllowed, "method_not_allowed")
	}
}

/*
handleAccount routes the request to the appropriate handler for /account endpoint.
*/
func (s *ApiServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
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
func (s *ApiServer) handleUniqueAccount(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetAccount(w, r)
	case "DELETE":
		return s.handleDeleteAccount(w, r)
	default:
		return NewApiError(http.StatusMethodNotAllowed, "method_not_allowed")
	}
}

func (s *ApiServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		return s.handleCreateTransaction(w, r)
	default:
		return NewApiError(http.StatusMethodNotAllowed, "method_not_allowed")
	}
}

/* ------------------------------- Controllers ------------------------------ */

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

/*
handleCreateToken is the controller that handles the POST /auth/login endpoint.
It creates a new token for the user.
*/
func (s *ApiServer) handleCreateToken(w http.ResponseWriter, r *http.Request) error {
	data := new(LoginDTO)

	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return NewApiError(http.StatusBadRequest, "invalid_request_body")
	}
	defer r.Body.Close()

	if data.AccountNumber <= 0 {
		return NewApiError(http.StatusBadRequest, "missing_account_number")
	}

	// Check if the account exists
	a, err := s.store.GetAccount(data.AccountNumber)
	if err != nil {
		if err.Error() == "account_not_found" {
			return NewApiError(http.StatusNotFound, err.Error())
		}
		return err
	}

	// TODO Check if the password is correct
	token, err := createAuthToken(a)

	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, NewApiResponse(http.StatusOK, token, r))
}

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
	id, err := getIntParameter(r, "id")
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
	id, err := getIntParameter(r, "id")
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

/*
	WIP

handleTransfer is the controller that handles the POST /transfer endpoint.
*/
func (s *ApiServer) handleCreateTransaction(w http.ResponseWriter, r *http.Request) error {
	data := new(CreateTransactionDTO)

	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return NewApiError(http.StatusBadRequest, "invalid_request_body")
	}
	defer r.Body.Close()

	if data.Amount <= 0 {
		return NewApiError(http.StatusBadRequest, "invalid_amount")
	} else if data.To == 0 {
		return NewApiError(http.StatusBadRequest, "invalid_to_account_id")
	}

	id := GetAuthenticatedAccountId(r)

	fromAccount, err := s.store.GetAccount(*id)
	if err != nil {
		if err.Error() == "account_not_found" {
			return NewApiError(http.StatusNotFound, "from_account_not_found")
		}
		return err
	}

	// Check if the to account exists
	toAccount, err := s.store.GetAccount(data.To)
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

	// Create the transaction

	err = s.store.CreateTxn(fromAccount.ID, data)
	if err != nil {
		return err
	}

	// Update the balance of the from account
	//fromBalance := balance - data.Amount

	// Update the balance of the to account
	//toBalance := balance + data.Amount

	return WriteJSON(w, http.StatusCreated, NewApiResponse(http.StatusCreated, data, r))
}
