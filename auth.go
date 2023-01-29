package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/farischt/gobank/config"
	jwt "github.com/golang-jwt/jwt/v4"
)

/*
withAuth is a middleware to protect routes that require authentication.
*/
func withAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("authentication protected route")

		token := r.Header.Get(config.GetConfig().GetString(config.TOKEN_NAME))
		if len(token) == 0 {
			WriteJSON(w, http.StatusUnauthorized, NewApiError(http.StatusUnauthorized, "missing_token"))
			return
		}

		t, err := validateAuthToken(token)
		if err != nil {
			WriteJSON(w, http.StatusUnauthorized, NewApiError(http.StatusUnauthorized, "invalid_token_error"))
			return
		}

		if !t.Valid {
			WriteJSON(w, http.StatusUnauthorized, NewApiError(http.StatusUnauthorized, "invalid_token"))
			return
		}

		// Equivalent to next() in express
		handlerFunc(w, r)
	}
}

/*
withoutAuth is a middleware to protect routes that must not be authenticated.
*/
func withoutAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Println("without authentication protected route")

		// Check if the token is already set
		token := r.Header.Get(config.GetConfig().GetString(config.TOKEN_NAME))
		if len(token) > 0 {
			WriteJSON(w, http.StatusForbidden, NewApiError(http.StatusForbidden, "already_authenticated"))
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
func createAuthToken(account *Account) (string, error) {
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
