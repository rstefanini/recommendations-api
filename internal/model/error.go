package model

import (
	"encoding/json"
	"errors"
	"net/http"
)

var ErrInvalidUserID = errors.New("")

type APIError interface {
	Message() string
	Status() int
	Error() *error
	WriteJSONError(w http.ResponseWriter)
}

type error struct {
	ErrorMessage string `json:"message"`
	ErrorStatus  int    `json:"status"`
}

func (e *error) Message() string {
	return e.ErrorMessage
}

func (e *error) Status() int {
	return e.ErrorStatus
}

func (e *error) Error() *error {
	return e
}

func (e *error) WriteJSONError(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(e.Status())

	json.NewEncoder(w).Encode(e)
}

func NewApiError(message string, status int) APIError {
	return &error{ErrorMessage: message, ErrorStatus: status}
}

func NewNotFoundApiError(message string) APIError {
	return NewApiError(message, http.StatusNotFound)
}

func NewBadRequestApiError(message string) APIError {
	return NewApiError(message, http.StatusBadRequest)
}

func NewInternalServerApiError(message string) APIError {
	return NewApiError(message, http.StatusInternalServerError)
}
