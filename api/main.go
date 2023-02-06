package api

import (
	"log"
	"net/http"

	"github.com/farischt/gobank/store"
	"github.com/gorilla/mux"
)

type Handlers struct {
	User *UserHandler
	Account *AccountHandler
	Transaction *TransactionHandler
	Authentication *AuthenticationHandler
}

func NewHandlers(store store.Store) *Handlers {
	return &Handlers{
		User: NewUserHandler(store),
		Account: NewAccountHandler(store),
		Transaction: NewTransactionHandler(store),
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
	 router.HandleFunc("/auth/login", WithoutAuth(makeHTTPFunc(s.handlers.Authentication.HandleLogin)))
	 router.HandleFunc("/account", makeHTTPFunc(s.handlers.Account.HandleAccount))
	 router.HandleFunc("/account/{id}", makeHTTPFunc(s.handlers.Account.HandleUniqueAccount))
	 router.HandleFunc("/transfer", WithAuth(makeHTTPFunc(s.handlers.Transaction.HandleTransfer)))

	log.Println("Server up and running on port", s.listenAddr[1:])
	err := http.ListenAndServe(s.listenAddr, router)

	if err != nil {
		log.Fatal(err)
	}
}

