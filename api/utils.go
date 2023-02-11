package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/farischt/gobank/config"
	"github.com/gorilla/mux"
)

/* API ERROR */
type ApiError struct {
	Status    int       `json:"status"`
	Err       string    `json:"error"`
	Timestamp time.Time `json:"timestamp"`
}

func NewApiError(s int, e string) ApiError {
	return ApiError{
		Status:    s,
		Err:       e,
		Timestamp: time.Now(),
	}
}

func (a ApiError) Error() string {
	return a.Err
}

/* API RESPONSE */
type ApiResponse struct {
	Status    int         `json:"status"`
	Timestamp time.Time   `json:"timestamp"`
	Method    string      `json:"method"`
	Path      string      `json:"path"`
	Data      interface{} `json:"data"`
}

func NewApiResponse(s int, d interface{}, r *http.Request) ApiResponse {
	return ApiResponse{
		Status:    s,
		Timestamp: time.Now(),
		Data:      d,
		Method:    r.Method,
		Path:      r.URL.Path,
	}
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
				_ = WriteJSON(w, e.Status, e)
				return
			}
			_ = WriteJSON(w, http.StatusInternalServerError, NewApiError(http.StatusInternalServerError, err.Error()))
		}
	}
}

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
getStringParameter is a helper function to get a string parameter from the request.
It takes the request and the parameter name.
It returns the parameter value and an error if the parameter is missing.
*/
func GetStringParameter(r *http.Request, param string) (string, error) {
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
func GetIntParameter(r *http.Request, param string) (uint, error) {
	p, err := GetStringParameter(r, param)
	if err != nil {
		return 0, err
	}

	parsedParameter, err := strconv.Atoi(p)
	if err != nil {
		return 0, NewApiError(http.StatusBadRequest, fmt.Sprintf("invalid_%s", param))
	}

	return uint(parsedParameter), nil
}

func GetTokenFromHeader(r *http.Request) (string, error) {
	token := r.Header.Get(config.GetConfig().GetString(config.TOKEN_NAME))
	if token == "" {
		return "", NewApiError(http.StatusUnauthorized, "missing_token")
	}

	return token, nil
}
