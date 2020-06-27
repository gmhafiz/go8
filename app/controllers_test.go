package app

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/matryer/is"
)

func TestHandleIndex(t *testing.T) {
	isTest := is.New(t)
	srv := Server{
		db: nil,
	}
	r := srv.RegisterRoutes()
	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	isTest.Equal(w.Code, http.StatusOK)
}
