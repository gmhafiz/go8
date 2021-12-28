package http

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/golang/mock/gomock"
	"github.com/jinzhu/now"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"

	"github.com/gmhafiz/go8/internal/domain/book"
	"github.com/gmhafiz/go8/internal/domain/book/mock"
	"github.com/gmhafiz/go8/internal/models"
)

//go:generate mockgen -package mock -source ../../handler.go -destination=../../mock/mock_handler.go

func TestHandler_Create(t *testing.T) {
	testBookRequest := &models.Book{
		Title:         "test01",
		PublishedDate: now.MustParse("2020-02-02"),
		ImageURL: null.String{
			String: "https://example.com/image.png",
			Valid:  true,
		},
		Description: "test01",
	}

	body, err := json.Marshal(testBookRequest)
	assert.NoError(t, err)
	ctrl := gomock.NewController(t)
	// If you are using a Go version of 1.14+, a mockgen version of 1.5.0+, and
	// are passing a *testing.T into gomock.NewController(t) you no longer need
	// to call ctrl.Finish() explicitly. It will be called for you automatically
	// from a self registered Cleanup function.
	defer ctrl.Finish()

	uc := mock.NewMockUseCase(ctrl)

	ctx := context.Background()
	var e error

	ucResp := &models.Book{
		ID:            1,
		Title:         testBookRequest.Title,
		PublishedDate: testBookRequest.PublishedDate,
		ImageURL:      testBookRequest.ImageURL,
		Description:   testBookRequest.Description,
	}
	uc.EXPECT().Create(ctx, testBookRequest).Return(ucResp, e).Times(1)

	router := chi.NewRouter()

	val := validator.New()
	h := RegisterHTTPEndPoints(router, val, uc)

	ww := httptest.NewRecorder()
	rr := httptest.NewRequest(http.MethodPost, "/api/v1/books", bytes.NewBuffer(body))

	h.Create(ww, rr)

	var gotBook book.Res
	err = json.NewDecoder(ww.Body).Decode(&gotBook)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, http.StatusCreated, ww.Code)
	assert.Equal(t, gotBook.Title, ucResp.Title)
	assert.Equal(t, gotBook.Description.String, ucResp.Description)
	assert.Equal(t, gotBook.PublishedDate.String(), ucResp.PublishedDate.String())
	assert.Equal(t, gotBook.ImageURL.String, ucResp.ImageURL.String)
}

func TestHandler_Get(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc := mock.NewMockUseCase(ctrl)

	ctx := context.Background()
	var e error
	var id int64
	ucResp := &models.Book{
		ID:            0,
		Title:         "",
		PublishedDate: time.Time{},
		ImageURL:      null.String{},
		Description:   "",
		CreatedAt:     null.Time{},
		UpdatedAt:     null.Time{},
		DeletedAt:     null.Time{},
	}

	uc.EXPECT().Read(ctx, id).Return(ucResp, e).Times(1)

	router := chi.NewRouter()

	h := NewHandler(uc, nil)
	RegisterHTTPEndPoints(router, nil, uc)

	ww := httptest.NewRecorder()
	rr, err := http.NewRequest(http.MethodGet, "/api/v1/books/1", nil)
	assert.NoError(t, err)

	h.Get(ww, rr)

	var gotBook book.Res
	err = json.NewDecoder(ww.Body).Decode(&gotBook)

	assert.NoError(t, err)
	assert.Equal(t, ucResp.ID, gotBook.ID)
	assert.Equal(t, ucResp.Description, gotBook.Description.String)
	assert.Equal(t, ucResp.PublishedDate, gotBook.PublishedDate)
	assert.Equal(t, ucResp.ImageURL, gotBook.ImageURL)
}

func TestHandler_List(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc := mock.NewMockUseCase(ctrl)

	ctx := context.Background()
	var e error

	uri := "/api/v1/books?page=1&size=30"
	f := book.Filters(url.Values{
		"page": []string{"1"},
		"size": []string{"30"},
	})

	var books []*models.Book

	uc.EXPECT().List(ctx, f).Return(books, e).AnyTimes()

	router := chi.NewRouter()

	val := validator.New()
	h := NewHandler(uc, val)
	RegisterHTTPEndPoints(router, val, uc)

	ww := httptest.NewRecorder()
	rr, err := http.NewRequest(http.MethodGet, uri, nil)
	assert.NoError(t, err)

	h.List(ww, rr)

	var gotBook []book.Res
	err = json.NewDecoder(ww.Body).Decode(&gotBook)

	assert.NoError(t, err)
}

func TestHandler_Update(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc := mock.NewMockUseCase(ctrl)

	ctx := context.Background()
	var e error
	bookReq := &models.Book{
		Title:         "test01",
		PublishedDate: now.MustParse("2020-02-02"),
		ImageURL: null.String{
			String: "https://example.com/image.png",
			Valid:  true,
		},
		Description: "test01",
	}
	body, err := json.Marshal(bookReq)
	assert.NoError(t, err)

	expectBook := &models.Book{
		ID:            1,
		Title:         "test01",
		PublishedDate: now.MustParse("2020-02-02"),
		ImageURL: null.String{
			String: "https://example.com/image.png",
			Valid:  true,
		},
		Description: "test01",
	}

	uc.EXPECT().Update(ctx, bookReq).Return(expectBook, e).Times(1)

	router := chi.NewRouter()

	val := validator.New()
	h := NewHandler(uc, val)
	RegisterHTTPEndPoints(router, val, uc)

	ww := httptest.NewRecorder()
	rr, err := http.NewRequest(http.MethodGet, "/api/v1/books/1", bytes.NewBuffer(body))
	assert.NoError(t, err)

	h.Update(ww, rr)

	var gotBook book.Res
	err = json.NewDecoder(ww.Body).Decode(&gotBook)

	assert.NoError(t, err)
}

func TestHandler_Delete(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	uc := mock.NewMockUseCase(ctrl)

	ctx := context.Background()
	var e error
	var id int64

	uc.EXPECT().Delete(ctx, id).Return(e).Times(1)

	router := chi.NewRouter()

	val := validator.New()
	h := NewHandler(uc, val)
	RegisterHTTPEndPoints(router, val, uc)

	ww := httptest.NewRecorder()
	rr, err := http.NewRequest(http.MethodGet, "/api/v1/books/1", nil)
	assert.NoError(t, err)

	h.Delete(ww, rr)

	assert.NoError(t, err)
}
