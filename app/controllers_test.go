package app

import (
	"eight/config"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
)

func TestHandleIndex(t *testing.T) {
	isTest := is.New(t)
	c := config.AppConfig()
	c.Testing = true
	server := NewApp(c)
	router := server.Router()

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	isTest.Equal(w.Code, http.StatusOK)
}
