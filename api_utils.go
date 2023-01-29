package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

/*
WriteJSON is a helper function to write JSON response.
It will set the content-type to application/json and write the status code.
It returns an error if the encoding fails.
*/
func WriteJSON(w http.ResponseWriter, status int, v any) error {
	w.Header().Set("content-type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

/*
apiFunc is a function that handles an API request.
It returns an error if the request fails.
*/
type apiFunc func(http.ResponseWriter, *http.Request) error

/*
makeHTTPFunc is a helper function to convert an apiFunc to http.HandlerFunc.
It returns an http.HandlerFunc that will write the error as JSON response.
*/
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

/*
getStringParameter is a helper function to get a string parameter from the request.
It takes the request and the parameter name.
It returns the parameter value and an error if the parameter is missing.
*/
func getStringParameter(r *http.Request, param string) (string, error) {
	vars := mux.Vars(r)
	p, ok := vars[param]
	if !ok {
		return "", NewApiError(http.StatusBadRequest, fmt.Sprintf("missing_%s", param))
	}

	return p, nil
}

/*
getIntParameter is a helper function to get an integer parameter from the request.
It takes the request and the parameter name.
It returns the parameter value and an error if the parameter is missing or invalid.
*/
func getIntParameter(r *http.Request, param string) (uint, error) {
	p, err := getStringParameter(r, param)
	if err != nil {
		return 0, err
	}

	parsedParameter, err := strconv.Atoi(p)
	if err != nil {
		return 0, NewApiError(http.StatusBadRequest, fmt.Sprintf("invalid_%s", param))
	}

	return uint(parsedParameter), nil
}
