package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotFoundHandler(t *testing.T) {
	server := Server{}

	req := httptest.NewRequest("GET", "/not-found", nil)
	rr := httptest.NewRecorder()

	server.notFoundHandler(rr, req)
	res := rr.Result()
	resBody, _ := io.ReadAll(res.Body)

	assert.Equal(t, http.StatusNotFound, res.StatusCode)
	assert.Empty(t, string(resBody))
}

func TestHealthHandler(t *testing.T) {
	server := Server{}

	req := httptest.NewRequest("GET", "/health", nil)
	rr := httptest.NewRecorder()

	server.healthHandler(rr, req)
	res := rr.Result()

	expected := "{\"health\":\"Healthy\"}"
	resBody, _ := io.ReadAll(res.Body)

	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.JSONEq(t, expected, string(resBody))
}
