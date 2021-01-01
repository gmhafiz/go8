package http

import (
	"bytes"
	"encoding/json"
	"github.com/gmhafiz/go8/internal/domain/book/usecase/mock"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/stretchr/testify/assert"
)

func TestHandler_Create(t *testing.T) {
	testBookRequest := &BookRequest{
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

	assert.Equal(t, 400, w.Code)
}
