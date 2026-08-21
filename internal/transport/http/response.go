package httptransport

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    any         `json:"data,omitempty"`
}

func WriteJSON(w http.ResponseWriter, status int, code int, message string, data any) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(Response{Code: code, Message: message, Data: data})
}

func OK(w http.ResponseWriter, data any) {
	WriteJSON(w, http.StatusOK, 0, "success", data)
}

func Fail(w http.ResponseWriter, status int, code int, message string) {
	WriteJSON(w, status, code, message, nil)
}

type APIError struct {
	Status  int
	Code    int
	Message string
}

func (e *APIError) Error() string { return e.Message }

func NewError(status, code int, message string) *APIError {
	return &APIError{Status: status, Code: code, Message: message}
}
