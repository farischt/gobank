package api

import (
	"log"
	"net/http"

	"github.com/farischt/gobank/config"
	"github.com/farischt/gobank/services"
	"github.com/farischt/gobank/store"
	"github.com/gorilla/mux"
)

type Handlers struct {
	User           *UserHandler
	Account        *AccountHandler
	Transaction    *TransactionHandler
	Authentication *AuthenticationHandler
}

func NewHandlers(service *services.Service) *Handlers {
	return &Handlers{
		User:           NewUserHandler(service),
		Account:        NewAccountHandler(service),
		Transaction:    NewTransactionHandler(service),
		Authentication: NewAuthenticationHandler(service),
	}
}

/*
ApiServer is the API server.
*/
type ApiServer struct {
	listenAddr string
	service    *services.Service
	handlers   *Handlers
}

/*
NewApiServer creates a new instance of API server.
*/
func New(l string, s store.Store) *ApiServer {
	services := services.New(s)

	return &ApiServer{
		listenAddr: l,
		service:    services,
		handlers:   NewHandlers(services),
	}
}

/*
Start starts the API server.
*/
func (s *ApiServer) Start() {
	router := mux.NewRouter()

	router.HandleFunc("/user", makeHTTPFunc(s.handlers.User.HandleUser))
	router.HandleFunc("/user/{id}", makeHTTPFunc(s.handlers.User.HandleUniqueUser))
	router.HandleFunc("/auth/login", s.WithoutAuth(makeHTTPFunc(s.handlers.Authentication.HandleLogin)))
	router.HandleFunc("/account", makeHTTPFunc(s.handlers.Account.HandleAccount))
	router.HandleFunc("/account/{id}", makeHTTPFunc(s.handlers.Account.HandleUniqueAccount))
	router.HandleFunc("/transfer", s.WithAuth(makeHTTPFunc(s.handlers.Transaction.HandleTransfer)))

	log.Println("Server up and running on port", s.listenAddr[1:])
	err := http.ListenAndServe(s.listenAddr, router)

	if err != nil {
		log.Fatal(err)
	}
}

/*
withAuth is a middleware to protect routes that require authentication.
*/
func (s *ApiServer) WithAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("authentication protected route")

		token := r.Header.Get(config.GetConfig().GetString(config.TOKEN_NAME))
		if len(token) == 0 {
			_ = WriteJSON(w, http.StatusUnauthorized, NewApiError(http.StatusUnauthorized, "missing_token"))
			return
		}

		_, validToken := s.service.Session.IsValidSessionToken(token)

		if !validToken {
			_ = WriteJSON(w, http.StatusUnauthorized, NewApiError(http.StatusUnauthorized, "invalid_token"))
			return
		}

		// Equivalent to next() in express
		handlerFunc(w, r)
	}
}

/*
withoutAuth is a middleware to protect routes that must not be authenticated.
*/
func (s *ApiServer) WithoutAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("without authentication protected route")

		// Check if the token is already set
		token := r.Header.Get(config.GetConfig().GetString(config.TOKEN_NAME))
		if len(token) > 0 {
			_ = WriteJSON(w, http.StatusForbidden, NewApiError(http.StatusForbidden, "already_authenticated"))
			return
		}

		// Equivalent to next() in express
		handlerFunc(w, r)
	}
}
