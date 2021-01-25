package http

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"

	"github.com/gmhafiz/go8/internal/domain/book"
	"github.com/gmhafiz/go8/internal/mock"
)

func TestHandler_Create(t *testing.T) {
	testBookRequest := &book.Request{
		Title:       "test01",
		Description: "test01",
	}

	r := chi.NewRouter()

	uc := new(mock.BookUseCaseMock)

	RegisterHTTPEndPoints(r, uc)

	body, err := json.Marshal(testBookRequest)
	assert.NoError(t, err)

	uc.On("Create", testBookRequest.Title, testBookRequest.Description).Return(nil)

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodPost, "/api/v1/books", bytes.NewBuffer(body))
	r.ServeHTTP(w, req)

	assert.Equal(t, http.StatusBadRequest, w.Code)
}
