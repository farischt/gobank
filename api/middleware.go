package api

import (
	"log"
	"net/http"

	"github.com/farischt/gobank/config"
)

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