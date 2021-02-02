package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/gmhafiz/go8/internal/domain/book"
	"github.com/gmhafiz/go8/internal/domain/book/usecase/mock"
	"github.com/gmhafiz/go8/internal/models"
)

//go:generate mockgen -package mock -source handler.go -aux_files handler=datastore.go -destination=handler_test.go

func TestHandler_Create(t *testing.T) {
	testBookRequest := &book.Request{
		Title:         "test01",
		PublishedDate: "2020-02-02",
		ImageURL:      "http://example.com/image.png",
		Description:   "test01",
	}
	body, err := json.Marshal(testBookRequest)
	assert.NoError(t, err)
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc := mock.NewMockUseCase(ctrl)

	ctx := context.Background()
	var e error

	var ucResp *models.Book
	uc.EXPECT().Create(ctx, *testBookRequest).Return(ucResp, e).AnyTimes()

	router := chi.NewRouter()

	h := NewHandler(uc)
	RegisterHTTPEndPoints(router, uc)

	ww := httptest.NewRecorder()
	rr, _ := http.NewRequest(http.MethodPost, "/api/v1/books", bytes.NewBuffer(body))

	h.Create(ww, rr)

	var gotBook book.Res
	err = json.NewDecoder(ww.Body).Decode(&gotBook)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusCreated, ww.Code)
}
