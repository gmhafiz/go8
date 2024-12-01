package handler

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"

	"github.com/gmhafiz/go8/internal/domain/book"
	"github.com/gmhafiz/go8/internal/domain/book/usecase"
	"github.com/gmhafiz/go8/internal/utility/message"
)

type Errs struct {
	Message []string `json:"message"`
}

func TestHandler_Create(t *testing.T) {
	type invalidCreateRequest struct {
		PublishedDate string `json:"published_date"`
	}

	type args struct {
		*book.CreateRequest
		invalidCreateRequest

		router    *chi.Mux
		validator *validator.Validate
	}
	type want struct {
		usecase struct {
			book *book.Schema
			err  error
		}
		res    *book.Res
		err    error
		errs   Errs
		status int
	}

	parsedTime, err := time.Parse(time.RFC3339, "2022-03-07T00:00:00Z")
	assert.NoError(t, err)

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "simple",
			args: args{
				CreateRequest: &book.CreateRequest{
					Title:         "Test Title",
					PublishedDate: "2022-03-07T00:00:00Z",
					ImageURL:      "https://example.com/image-test.png",
					Description:   "Test Description",
				},
				router:    chi.NewRouter(),
				validator: validator.New(),
			},
			want: want{
				usecase: struct {
					book *book.Schema
					err  error
				}{
					book: &book.Schema{
						ID:            1,
						Title:         "Test Title",
						PublishedDate: parsedTime,
						ImageURL:      "https://example.com/image-test.png",
						Description:   "Test Description",
						CreatedAt:     time.Now(),
						UpdatedAt:     time.Now(),
						DeletedAt:     sql.NullTime{},
					},
					err: nil,
				},
				res: &book.Res{
					ID:            1,
					Title:         "Test Title",
					PublishedDate: parsedTime,
					ImageURL:      "https://example.com/image-test.png",
					Description:   "Test Description",
				},
				status: http.StatusCreated,
			},
		},
		{
			name: "invalid request",
			args: args{
				invalidCreateRequest: invalidCreateRequest{
					PublishedDate: "2022-03-07T00:00:00Z",
				},
				router:    chi.NewRouter(),
				validator: validator.New(),
			},
			want: want{
				usecase: struct {
					book *book.Schema
					err  error
				}{
					book: &book.Schema{},
					err:  nil,
				},
				res: &book.Res{},
				errs: Errs{Message: []string{
					"Title is required with type string",
					"ImageURL is url with type string",
					"Description is required with type string",
				}},
				status: http.StatusBadRequest,
			},
		},
		{
			name: "no row",
			args: args{
				CreateRequest: &book.CreateRequest{
					Title:         "Title",
					PublishedDate: "2022-03-07T00:00:00Z",
					ImageURL:      "https://example.com/image-test.png",
					Description:   "Description",
				},
				router:    chi.NewRouter(),
				validator: validator.New(),
			},
			want: want{
				usecase: struct {
					book *book.Schema
					err  error
				}{
					book: &book.Schema{},
					err:  sql.ErrNoRows,
				},
				res:    &book.Res{},
				err:    message.ErrBadRequest,
				status: http.StatusBadRequest,
			},
		},
		{
			name: "other error",
			args: args{
				CreateRequest: &book.CreateRequest{
					Title:         "Title",
					PublishedDate: "2022-03-07T00:00:00Z",
					ImageURL:      "https://example.com/image-test.png",
					Description:   "Description",
				},
				router:    chi.NewRouter(),
				validator: validator.New(),
			},
			want: want{
				usecase: struct {
					book *book.Schema
					err  error
				}{
					book: &book.Schema{},
					err:  errors.New("other error"),
				},
				res:    &book.Res{},
				err:    errors.New("other error"),
				status: http.StatusInternalServerError,
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer
			var err error
			if tt.args.CreateRequest != nil {
				err = json.NewEncoder(&buf).Encode(tt.args.CreateRequest)
			} else {
				err = json.NewEncoder(&buf).Encode(tt.args.invalidCreateRequest)
			}
			assert.Nil(t, err)

			rr := httptest.NewRequest(http.MethodPost, "/api/v1/book", &buf)
			ww := httptest.NewRecorder()

			uc := &usecase.BookMock{
				CreateFunc: func(ctx context.Context, bookMiripParam *book.CreateRequest) (*book.Schema, error) {
					return tt.want.usecase.book, tt.want.usecase.err
				},
			}

			h := RegisterHTTPEndPoints(tt.args.router, tt.args.validator, uc)

			h.Create(ww, rr)

			//if tt.args.CreateRequest == nil {
			//	var errs Errs
			//	if err := json.NewDecoder(ww.Body).Decode(&errs); err != nil {
			//		t.Fatal(err)
			//	}
			//	assert.Equal(t, tt.want.errs, errs)
			//} else {

			assert.Equal(t, tt.want.status, ww.Code)

			if ww.Code >= 200 && ww.Code < 300 {
				var got book.Res

				if err = json.NewDecoder(ww.Body).Decode(&got); err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, tt.want.res.ID, got.ID)
				assert.Equal(t, tt.want.res.Title, got.Title)
				assert.Equal(t, tt.want.res.PublishedDate, got.PublishedDate)
				assert.Equal(t, tt.want.res.ImageURL, got.ImageURL)
				assert.Equal(t, tt.want.res.Description, got.Description)
			} else {
				b, err := io.ReadAll(ww.Body)
				assert.Nil(t, err)

				if len(tt.want.errs.Message) > 0 {
					errStruct := Errs{}

					err = json.Unmarshal(b, &errStruct)
					assert.Nil(t, err)

					for i := range errStruct.Message {
						assert.Equal(t, tt.want.errs.Message[i], errStruct.Message[i])
					}

				} else {
					errStruct := struct {
						Message string `json:"message"`
					}{
						Message: string(b),
					}

					err = json.Unmarshal(b, &errStruct)
					assert.Nil(t, err)
					assert.Equal(t, tt.want.err.Error(), errStruct.Message)
				}
			}

			//}
		})
	}
}

func TestHandler_Get(t *testing.T) {
	type args struct {
		bookID int
		param  string

		router    *chi.Mux
		validator *validator.Validate
	}
	type want struct {
		usecase struct {
			book *book.Schema
			err  error
		}
		res    *book.Res
		err    error
		status int
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "simple",
			args: args{
				bookID:    1,
				param:     "bookID",
				router:    chi.NewRouter(),
				validator: validator.New(),
			},
			want: want{
				usecase: struct {
					book *book.Schema
					err  error
				}{
					&book.Schema{
						ID:            1,
						Title:         "",
						PublishedDate: time.Time{},
						ImageURL:      "",
						Description:   "",
					},
					nil,
				},
				status: http.StatusOK,
			},
		},
		{
			name: "wrong URL param",
			args: args{
				bookID:    1,
				param:     "id",
				router:    chi.NewRouter(),
				validator: validator.New(),
			},
			want: want{
				usecase: struct {
					book *book.Schema
					err  error
				}{
					&book.Schema{},
					nil,
				},
				res:    &book.Res{},
				err:    message.ErrBadRequest,
				status: http.StatusBadRequest,
			},
		},
		{
			name: "no record found",
			args: args{
				bookID:    1,
				param:     "bookID",
				router:    chi.NewRouter(),
				validator: validator.New(),
			},
			want: want{
				usecase: struct {
					book *book.Schema
					err  error
				}{
					&book.Schema{},
					sql.ErrNoRows,
				},
				res:    &book.Res{},
				err:    errors.New("no book is found for this ID"),
				status: http.StatusBadRequest,
			},
		},
		{
			name: "other error",
			args: args{
				bookID:    1,
				param:     "bookID",
				router:    chi.NewRouter(),
				validator: validator.New(),
			},
			want: want{
				usecase: struct {
					book *book.Schema
					err  error
				}{
					&book.Schema{},
					errors.New("some other error"),
				},
				res:    &book.Res{},
				err:    nil,
				status: http.StatusInternalServerError,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			rr := httptest.NewRequest(http.MethodGet, fmt.Sprintf("/api/v1/book/{%s}", tt.args.param), nil)
			ww := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add(tt.args.param, strconv.Itoa(tt.args.bookID))

			rr = rr.WithContext(context.WithValue(rr.Context(), chi.RouteCtxKey, rctx))

			uc := &usecase.BookMock{
				ReadFunc: func(ctx context.Context, bookID uint64) (*book.Schema, error) {
					return tt.want.usecase.book, tt.want.usecase.err
				},
			}

			h := RegisterHTTPEndPoints(tt.args.router, tt.args.validator, uc)

			h.Get(ww, rr)

			assert.Equal(t, tt.want.status, ww.Code)

			if ww.Code >= 200 && ww.Code < 300 {
				var got book.Schema
				if err := json.NewDecoder(ww.Body).Decode(&got); err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, tt.want.usecase.book.ID, got.ID)
				assert.Equal(t, tt.want.usecase.book.Title, got.Title)
				assert.Equal(t, tt.want.usecase.book.PublishedDate, got.PublishedDate)
				assert.Equal(t, tt.want.usecase.book.ImageURL, got.ImageURL)
				assert.Equal(t, tt.want.usecase.book.Description, got.Description)
			} else {
				b, err := io.ReadAll(ww.Body)
				assert.Nil(t, err)

				errStruct := struct {
					Message string `json:"message"`
				}{
					Message: string(b),
				}

				if len(b) == 0 {
					return
				}

				err = json.Unmarshal(b, &errStruct)
				assert.Nil(t, err)
				assert.Equal(t, tt.want.err, errors.New(errStruct.Message))
			}
		})
	}
}

func TestHandler_List(t *testing.T) {
	type args struct {
	}
	type want struct {
		usecase struct {
			books []*book.Schema
			err   error
		}

		books  []*book.Res
		err    error
		status int
	}

	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "simple",
			args: args{},
			want: want{
				usecase: struct {
					books []*book.Schema
					err   error
				}{
					books: []*book.Schema{
						{
							ID:            1,
							Title:         "test title",
							PublishedDate: time.Time{},
							ImageURL:      "",
							Description:   "",
							CreatedAt:     time.Time{},
							UpdatedAt:     time.Time{},
							DeletedAt:     sql.NullTime{},
						},
					},
					err: nil,
				},
				books: []*book.Res{
					{
						ID:            1,
						Title:         "test title",
						PublishedDate: time.Time{},
						ImageURL:      "",
						Description:   "",
					},
				},
				status: http.StatusOK,
			},
		},
		{
			name: "error fetching books",
			args: args{},
			want: want{
				usecase: struct {
					books []*book.Schema
					err   error
				}{
					books: nil,
					err:   message.ErrFetchingBook,
				},
				books:  []*book.Res{},
				err:    message.ErrFetchingBook,
				status: http.StatusInternalServerError,
			},
		},
		{
			name: "other errors (internal)",
			args: args{},
			want: want{
				usecase: struct {
					books []*book.Schema
					err   error
				}{
					books: nil,
					err:   errors.New("some other error"),
				},
				books:  []*book.Res{},
				err:    errors.New("some other error"),
				status: http.StatusInternalServerError,
			},
		},
		{
			name: "return no book",
			args: args{},
			want: want{
				usecase: struct {
					books []*book.Schema
					err   error
				}{
					books: []*book.Schema{},
					err:   nil,
				},
				books:  []*book.Res{},
				err:    nil,
				status: http.StatusOK,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			rr := httptest.NewRequest(http.MethodGet, "/api/v1/book", nil)
			ww := httptest.NewRecorder()

			uc := &usecase.BookMock{
				ListFunc: func(ctx context.Context, f *book.Filter) ([]*book.Schema, error) {
					return tt.want.usecase.books, tt.want.usecase.err
				},
			}

			h := RegisterHTTPEndPoints(chi.NewRouter(), validator.New(), uc)

			h.List(ww, rr)

			assert.Equal(t, tt.want.status, ww.Code)

			if ww.Code >= 200 && ww.Code < 300 {
				var got []*book.Res
				if err := json.NewDecoder(ww.Body).Decode(&got); err != nil {
					t.Fatal(err)
				}

				for i := 0; i < len(got); i++ {
					for j := 0; j < len(tt.want.books); j++ {
						assert.Equal(t, tt.want.books[j].ID, got[i].ID)
						assert.Equal(t, tt.want.books[j].Title, got[i].Title)
						assert.Equal(t, tt.want.books[j].Description, got[i].Description)
						assert.Equal(t, tt.want.books[j].ImageURL, got[i].ImageURL)
						assert.Equal(t, tt.want.books[j].PublishedDate.String(), got[i].PublishedDate.String())
					}
				}

				return
			} else {
				b, err := io.ReadAll(ww.Body)
				assert.Nil(t, err)

				errStruct := struct {
					Message string `json:"message"`
				}{
					Message: string(b),
				}

				err = json.Unmarshal(b, &errStruct)
				assert.Nil(t, err)
				assert.Equal(t, tt.want.err.Error(), errStruct.Message)
			}

		})
	}
}

func TestHandler_Update(t *testing.T) {
	parsedTime, err := time.Parse(time.RFC3339, "2022-03-09T00:00:00Z")
	assert.Nil(t, err)

	type invalidRequest struct {
		Title string `json:"title"`
	}

	type args struct {
		bookID int
		book   book.UpdateRequest
		invalidRequest
		param string
	}
	type want struct {
		usecase struct {
			book *book.Schema
			err  error
		}
		status int
		book   *book.Res
		err    error
		errs   Errs
	}

	tests := []struct {
		name string
		args
		want
	}{
		{
			name: "simple",
			args: args{
				bookID: 1,
				book: book.UpdateRequest{
					ID:            1,
					Title:         "mock title",
					PublishedDate: "2022-03-09T00:00:00Z",
					ImageURL:      "https://example.com/image.png",
					Description:   "mock description",
				},
				param: "bookID",
			},
			want: want{
				usecase: struct {
					book *book.Schema
					err  error
				}{
					book: &book.Schema{
						ID:            1,
						Title:         "mock title",
						PublishedDate: parsedTime,
						ImageURL:      "https://example.com/image.png",
						Description:   "mock description",
					},
					err: nil,
				},
				book: &book.Res{
					ID:            1,
					Title:         "mock title",
					PublishedDate: parsedTime,
					ImageURL:      "https://example.com/image.png",
					Description:   "mock description",
				},
				status: http.StatusOK,
			},
		},
		{
			name: "invalid query parameter",
			args: args{
				bookID: 1,
				book: book.UpdateRequest{
					ID:            1,
					Title:         "mock title",
					PublishedDate: "2022-03-09T00:00:00Z",
					ImageURL:      "https://example.com/image.png",
					Description:   "mock description",
				},
				param: "id",
			},
			want: want{
				usecase: struct {
					book *book.Schema
					err  error
				}{},
				book:   &book.Res{},
				status: http.StatusBadRequest,
				err:    message.ErrBadRequest,
			},
		},
		{
			name: "insufficient update request payload",
			args: args{
				bookID: 1,
				invalidRequest: invalidRequest{
					Title: "only title",
				},
				param: "bookID",
			},
			want: want{
				usecase: struct {
					book *book.Schema
					err  error
				}{
					book: &book.Schema{},
					err:  nil,
				},
				status: http.StatusBadRequest,
				book:   &book.Res{},
				errs: Errs{Message: []string{
					"PublishedDate is required with type string",
					"ImageURL is url with type string",
					"Description is required with type string",
				}},
			},
		},
		{
			name: "some other internal error",
			args: args{
				bookID: 1,
				book: book.UpdateRequest{
					ID:            1,
					Title:         "mock title",
					PublishedDate: "2022-03-09T00:00:00Z",
					ImageURL:      "https://example.com/image.png",
					Description:   "mock description",
				},
				param: "bookID",
			},
			want: want{
				usecase: struct {
					book *book.Schema
					err  error
				}{
					&book.Schema{},
					errors.New("some error"),
				},
				status: http.StatusInternalServerError,
				book:   &book.Res{},
				err:    errors.New("some error"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			var buf bytes.Buffer
			if tt.args.invalidRequest.Title != "" {
				err := json.NewEncoder(&buf).Encode(tt.args.invalidRequest)
				assert.Nil(t, err)
			} else {
				err := json.NewEncoder(&buf).Encode(tt.args.book)
				assert.Nil(t, err)
			}

			rr := httptest.NewRequest(http.MethodPut, fmt.Sprintf("/api/v1/book/{%s}", tt.args.param), &buf)
			ww := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add(tt.args.param, strconv.Itoa(tt.args.bookID))

			rr = rr.WithContext(context.WithValue(rr.Context(), chi.RouteCtxKey, rctx))

			uc := &usecase.BookMock{
				UpdateFunc: func(ctx context.Context, bookParam *book.UpdateRequest) (*book.Schema, error) {
					return tt.want.usecase.book, tt.want.usecase.err
				},
			}

			h := RegisterHTTPEndPoints(chi.NewRouter(), validator.New(), uc)

			h.Update(ww, rr)

			assert.Equal(t, tt.want.status, ww.Code)

			if ww.Code >= 200 && ww.Code < 300 {
				var got book.Res
				if err := json.NewDecoder(ww.Body).Decode(&got); err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, tt.want.book.ID, got.ID)
				assert.Equal(t, tt.want.book.Title, got.Title)
				assert.Equal(t, tt.want.book.PublishedDate, got.PublishedDate)
				assert.Equal(t, tt.want.book.ImageURL, got.ImageURL)
				assert.Equal(t, tt.want.book.Description, got.Description)
				return

			} else {
				b, err := io.ReadAll(ww.Body)
				assert.Nil(t, err)

				if len(tt.want.errs.Message) > 0 {
					errStruct := Errs{}

					err = json.Unmarshal(b, &errStruct)
					assert.Nil(t, err)

					for i := range errStruct.Message {
						assert.Equal(t, tt.want.errs.Message[i], errStruct.Message[i])
					}

				} else {
					errStruct := struct {
						Message string `json:"message"`
					}{
						Message: string(b),
					}

					err = json.Unmarshal(b, &errStruct)
					assert.Nil(t, err)
					assert.Equal(t, tt.want.err.Error(), errStruct.Message)
				}

			}

		})
	}
}

func TestHandler_Delete(t *testing.T) {
	type args struct {
		bookID int
		param  string
	}
	type want struct {
		status int
		error
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "ok",
			args: args{
				bookID: 1,
				param:  "bookID",
			},
			want: want{
				status: http.StatusOK,
				error:  nil,
			},
		},
		{
			name: "wrong query param",
			args: args{
				bookID: 1,
				param:  "id",
			},
			want: want{
				status: http.StatusBadRequest,
				error:  nil,
			},
		},
		{
			name: "some internal error",
			args: args{
				bookID: 1,
				param:  "bookID",
			},
			want: want{
				status: http.StatusInternalServerError,
				error:  errors.New("some internal error"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			rr := httptest.NewRequest(http.MethodDelete, fmt.Sprintf("/api/v1/book/{%s}", tt.args.param), nil)
			ww := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add(tt.args.param, strconv.Itoa(tt.args.bookID))

			rr = rr.WithContext(context.WithValue(rr.Context(), chi.RouteCtxKey, rctx))

			router := chi.NewRouter()
			val := validator.New()

			uc := &usecase.BookMock{
				DeleteFunc: func(ctx context.Context, bookID uint64) error {
					return tt.want.error
				},
			}

			h := RegisterHTTPEndPoints(router, val, uc)

			h.Delete(ww, rr)

			assert.Equal(t, tt.want.status, ww.Code)
		})
	}
}
