package main

import (
	"net/http"
	"time"
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
