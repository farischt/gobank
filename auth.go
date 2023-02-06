package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/farischt/gobank/config"
	"github.com/farischt/gobank/dto"
	"github.com/farischt/gobank/types"
	jwt "github.com/golang-jwt/jwt/v4"
)

func (s *ApiServer) HandleLogin(w http.ResponseWriter, r *http.Request) error {
	switch r.Method {
	case "POST":
		return s.handleCreateToken(w, r)
	default:
		return NewApiError(http.StatusMethodNotAllowed, "method_not_allowed")
	}
}

/*
handleCreateToken is the controller that handles the POST /auth/login endpoint.
It creates a new token for the user.
*/
func (s *ApiServer) handleCreateToken(w http.ResponseWriter, r *http.Request) error {
	data := new(dto.LoginDTO)

	if err := json.NewDecoder(r.Body).Decode(data); err != nil {
		return NewApiError(http.StatusBadRequest, "invalid_request_body")
	}
	defer r.Body.Close()

	if data.AccountNumber <= 0 {
		return NewApiError(http.StatusBadRequest, "missing_account_number")
	}

	// Check if the account exists
	a, err := s.store.Account.GetAccount(data.AccountNumber)
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
withAuth is a middleware to protect routes that require authentication.
*/
func WithAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("authentication protected route")

		token := r.Header.Get(config.GetConfig().GetString(config.TOKEN_NAME))
		if len(token) == 0 {
			_ = WriteJSON(w, http.StatusUnauthorized, NewApiError(http.StatusUnauthorized, "missing_token"))
			return
		}

		t, err := validateAuthToken(token)
		if err != nil {
			_ = WriteJSON(w, http.StatusUnauthorized, NewApiError(http.StatusUnauthorized, "invalid_token_error"))
			return
		}

		if !t.Valid {
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
func WithoutAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
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
validateAuth is a function to validate the token.
*/
func validateAuthToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		secret := config.TOKEN_SECRET
		return []byte(secret), nil
	})
}

/*
createAuthToken is a function to create a token.
It takes an account and a user as input.
It returns the token and an error.
*/
func createAuthToken(account *types.Account) (string, error) {
	claims := jwt.MapClaims{
		"expires_at": 150000,
		"account_id": account.ID,
	}

	secret := config.TOKEN_SECRET
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

/*
decryptAuthToken is a function to decrypt the token.
*/
func decryptAuthToken(token *jwt.Token) map[string]interface{} {
	return token.Claims.(jwt.MapClaims)
}

/*
GetAuthenticatedAccountId is a function to get the authenticated account id from the jwt token.
*/
func GetAuthenticatedAccountId(r *http.Request) *uint {
	token := r.Header.Get(config.GetConfig().GetString(config.TOKEN_NAME))
	t, err := validateAuthToken(token)

	if err != nil || !t.Valid {
		return nil
	}

	claims := decryptAuthToken(t)
	accountId := uint(claims["account_id"].(float64))

	return &accountId
}
