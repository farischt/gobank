package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

/*
API SERVER
*/
type ApiServer struct {
	listenAddr string
	store      Storage
}

func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type apiFunc func(http.ResponseWriter, *http.Request) error

func makeHTTPFunc(f apiFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		err := f(w, r)
		if err != nil {
			if e, ok := err.(ApiError); ok {
				WriteJSON(w, e.Status, e)
				return
			}
			WriteJSON(w, http.StatusInternalServerError, NewApiError(http.StatusInternalServerError, err.Error()))
		}
	}
}

func NewApiServer(l string, s Storage) *ApiServer {
	return &ApiServer{
		listenAddr: l,
		store:      s,
	}
}

func (s *ApiServer) Start() {
	router := mux.NewRouter()

	router.HandleFunc("/account", makeHTTPFunc(s.handleAccount))
	router.HandleFunc("/account/{id}", makeHTTPFunc(s.handleGetAccount))

	log.Println("Server up and running on port ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

/*
ACCOUNT HANDLER
*/
func (s *ApiServer) handleAccount(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "GET":
		return s.handleGetAccounts(w, r)
	case "POST":
		return s.handleCreateAccount(w, r)
	case "DELETE":
		return s.handleDeleteAccount(w, r)
	default:
		return NewApiError(http.StatusMethodNotAllowed, "method_not_allowed")
	}
}

func (s *ApiServer) handleGetAccount(w http.ResponseWriter, r *http.Request) error {
	vars := mux.Vars(r)
	id, ok := vars["id"]
	if !ok {
		return NewApiError(http.StatusBadRequest, "missing_id")
	}

	parsedId, err := strconv.Atoi(id)
	if err != nil {
		return NewApiError(http.StatusBadRequest, "invalid_id")
	}

	a, err := s.store.GetAccount(parsedId)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			return NewApiError(http.StatusNotFound, "account_not_found")
		}
		return err
	}

	return WriteJSON(w, http.StatusOK, NewApiResponse(http.StatusOK, a, r))
}

func (s *ApiServer) handleGetAccounts(w http.ResponseWriter, r *http.Request) error {
	accounts, err := s.store.GetAllAccount()
	if err != nil {
		return err
	}

	return WriteJSON(w, http.StatusOK, NewApiResponse(http.StatusOK, accounts, r))
}

func (s *ApiServer) handleCreateAccount(w http.ResponseWriter, r *http.Request) error {
	data := new(CreateAccountDTO)

	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return NewApiError(http.StatusBadRequest, "invalid_request_body")
	}

	if len(data.FirstName) == 0 {
		return NewApiError(http.StatusBadRequest, "empty_first_name")
	} else if len(data.LastName) == 0 {
		return NewApiError(http.StatusBadRequest, "empty_last_name")
	}

	err := s.store.CreateAccount(data)
	if err != nil {
		return err
	}

	fmt.Println(err)

	return WriteJSON(w, http.StatusCreated, NewApiResponse(http.StatusCreated, data, r))
}

func (s *ApiServer) handleDeleteAccount(w http.ResponseWriter, r *http.Request) error {
	return nil
}

// func (s *ApiServer) handleTransfer(w http.ResponseWriter, r *http.Request) error {
// 	return nil
// }
