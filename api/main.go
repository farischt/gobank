package api

import (
	"log"
	"net/http"

	"github.com/farischt/gobank/config"
	"github.com/farischt/gobank/store"
	"github.com/gorilla/mux"
)

type Handlers struct {
	User           *UserHandler
	Account        *AccountHandler
	Transaction    *TransactionHandler
	Authentication *AuthenticationHandler
}

func NewHandlers(store store.Store) *Handlers {
	return &Handlers{
		User:           NewUserHandler(store),
		Account:        NewAccountHandler(store),
		Transaction:    NewTransactionHandler(store),
		Authentication: NewAuthenticationHandler(store),
	}
}

/*
ApiServer is the API server.
*/
type ApiServer struct {
	listenAddr string
	store      store.Store
	handlers   *Handlers
}

/*
NewApiServer creates a new instance of API server.
*/
func New(l string, s store.Store) *ApiServer {
	return &ApiServer{
		listenAddr: l,
		store:      s,
		handlers:   NewHandlers(s),
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

		_, validToken := s.store.SessionToken.IsValidSessionToken(token)

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

/*
GetAuthenticatedAccountId is a function to get the authenticated account id from the jwt token.
*/
func GetAuthenticatedAccountId(r *http.Request, s store.Store) *uint {
	token := r.Header.Get(config.GetConfig().GetString(config.TOKEN_NAME))

	t, ok := s.SessionToken.IsValidSessionToken(token)

	if ok {
		return &t
	}

	return nil
}
