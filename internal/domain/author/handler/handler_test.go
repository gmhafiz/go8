package handler

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"

	"github.com/gmhafiz/go8/internal/domain/author"
	"github.com/gmhafiz/go8/internal/domain/author/usecase"
	"github.com/gmhafiz/go8/internal/domain/book"
	"github.com/gmhafiz/go8/internal/utility/message"
	"github.com/gmhafiz/go8/internal/utility/respond"
)

var (
	ErrTransactionFailed = errors.New("simulate transaction failed")
)

type Errs struct {
	Message []string `json:"message"`
}

func TestHandler_Create(t *testing.T) {
	type invalidCreateRequest struct {
		LastName string `json:"last_name"`
	}
	type args struct {
		*author.CreateRequest
		invalidCreateRequest
	}

	type want struct {
		usecase struct {
			*author.Schema
			error
		}
		response *author.GetResponse
		err      error
		Errs
		status int
	}

	tests := []struct {
		name string
		args
		want
	}{
		{
			name: "simple",
			args: args{
				CreateRequest: &author.CreateRequest{
					FirstName:  "First",
					MiddleName: "Middle",
					LastName:   "Last",
					Books:      nil,
				},
			},
			want: want{
				usecase: struct {
					*author.Schema
					error
				}{
					&author.Schema{
						ID:         1,
						FirstName:  "First",
						MiddleName: "Middle",
						LastName:   "Last",
						CreatedAt:  time.Now(),
						UpdatedAt:  time.Now(),
						DeletedAt:  nil,
						Books:      make([]*book.Schema, 0),
					},
					nil,
				},
				response: &author.GetResponse{
					ID:         1,
					FirstName:  "First",
					MiddleName: "Middle",
					LastName:   "Last",
					Books:      make([]*book.Schema, 0),
				},
				err:    nil,
				status: http.StatusCreated,
			},
		},
		{
			name: "invalid create request",
			args: args{
				nil,
				invalidCreateRequest{
					LastName: "last Name",
				},
			},
			want: want{
				usecase: struct {
					*author.Schema
					error
				}{
					&author.Schema{},
					nil,
				},
				response: &author.GetResponse{},
				Errs: Errs{
					Message: []string{"FirstName is required with type string"},
				},
				status: http.StatusBadRequest,
			},
		},
		{
			name: "simulate transaction rollback",
			args: args{
				CreateRequest: &author.CreateRequest{
					FirstName:  "First",
					MiddleName: "Middle",
					LastName:   "Last",
					Books:      nil,
				},
			},
			want: want{
				usecase: struct {
					*author.Schema
					error
				}{
					&author.Schema{
						ID:         1,
						FirstName:  "First",
						MiddleName: "Middle",
						LastName:   "Last",
						CreatedAt:  time.Now(),
						UpdatedAt:  time.Now(),
						DeletedAt:  nil,
						Books:      make([]*book.Schema, 0),
					},
					ErrTransactionFailed,
				},
				response: &author.GetResponse{},
				err:      ErrTransactionFailed,
				status:   http.StatusInternalServerError,
			},
		},
		{
			name: "no row",
			args: args{
				CreateRequest: &author.CreateRequest{
					FirstName:  "First",
					MiddleName: "Middle",
					LastName:   "Last",
				},
			},
			want: want{
				usecase: struct {
					*author.Schema
					error
				}{
					&author.Schema{},
					sql.ErrNoRows,
				},
				response: &author.GetResponse{},
				err:      message.ErrBadRequest,
				status:   http.StatusBadRequest,
			},
		},
		{
			name: "other error",
			args: args{
				CreateRequest: &author.CreateRequest{
					FirstName:  "First",
					MiddleName: "Middle",
					LastName:   "Last",
				},
			},
			want: want{
				usecase: struct {
					*author.Schema
					error
				}{
					&author.Schema{},
					errors.New("other error"),
				},
				response: &author.GetResponse{},
				err:      errors.New("other error"),
				status:   http.StatusInternalServerError,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			var buf bytes.Buffer
			var err error
			if test.args.CreateRequest != nil {
				err = json.NewEncoder(&buf).Encode(test.args.CreateRequest)
			} else {
				err = json.NewEncoder(&buf).Encode(test.args.invalidCreateRequest)
			}
			assert.Nil(t, err)

			rr := httptest.NewRequest(http.MethodPost, "/api/v1/author", &buf)
			ww := httptest.NewRecorder()

			router := chi.NewRouter()
			val := validator.New()

			uc := &usecase.AuthorMock{
				CreateFunc: func(ctx context.Context, a *author.CreateRequest) (*author.Schema, error) {
					return test.want.usecase.Schema, test.want.usecase.error
				},
			}

			h := RegisterHTTPEndPoints(router, val, uc)
			h.Create(ww, rr)

			if test.args.CreateRequest == nil {
				var errs Errs
				if err = json.NewDecoder(ww.Body).Decode(&errs); err != nil {
					t.Fatal(err)
				}
				assert.Equal(t, test.want.Errs, errs)
			} else {
				assert.Equal(t, ww.Code, test.status)

				if ww.Code >= 200 && ww.Code < 300 {
					var got author.GetResponse
					if err = json.NewDecoder(ww.Body).Decode(&got); err != nil {
						t.Fatal(err)
					}

					assert.Equal(t, &got, test.want.response)
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
					assert.Equal(t, test.want.err.Error(), errStruct.Message)
				}
			}
		})
	}
}

func TestHandler_List(t *testing.T) {
	type args struct {
		uri string
	}

	type want struct {
		usecase struct {
			authors []*author.Schema
			total   int
			error
		}
		status int
		size   int
		total  int
		error
	}

	type test struct {
		name string
		args
		want
	}

	tests := []test{
		{
			name: "returns no record",
			args: args{
				uri: "/api/v1/author",
			},
			want: want{
				usecase: struct {
					authors []*author.Schema
					total   int
					error
				}{
					authors: make([]*author.Schema, 0),
					total:   0,
					error:   nil,
				},
				status: http.StatusOK,
				error:  nil,
				size:   0,
				total:  0,
			},
		},
		{
			name: "returns one record",
			args: args{
				uri: "/api/v1/author",
			},
			want: want{
				usecase: struct {
					authors []*author.Schema
					total   int
					error
				}{
					authors: []*author.Schema{
						{
							ID:         1,
							FirstName:  "First",
							MiddleName: "Middle",
							LastName:   "Last",
							CreatedAt:  time.Time{},
							UpdatedAt:  time.Time{},
							DeletedAt:  nil,
							Books:      nil,
						},
					},
					total: 1,
					error: nil,
				},
				status: http.StatusOK,
				error:  nil,
				size:   1,
				total:  1,
			},
		},
		{
			name: "simulate lower layer error",
			args: args{
				uri: "/api/v1/author",
			},
			want: want{
				status: http.StatusInternalServerError,
				size:   0,
				total:  0,
				usecase: struct {
					authors []*author.Schema
					total   int
					error
				}{
					authors: []*author.Schema{},
					total:   0,
					error:   errors.New("some use case internal error"),
				},
				error: errors.New("some use case internal error"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			rr := httptest.NewRequest(http.MethodGet, test.args.uri, nil)
			ww := httptest.NewRecorder()

			router := chi.NewRouter()

			val := validator.New()

			uc := &usecase.AuthorMock{
				ListFunc: func(ctx context.Context, f *author.Filter) ([]*author.Schema, int, error) {
					return test.want.usecase.authors, test.want.usecase.total, test.want.usecase.error
				},
			}

			h := RegisterHTTPEndPoints(router, val, uc)
			h.List(ww, rr)

			assert.Equal(t, test.want.status, ww.Code)

			if ww.Code >= 200 && ww.Code < 300 {
				var got respond.Standard
				if err := json.NewDecoder(ww.Body).Decode(&got); err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, test.want.size, got.Meta.Size)
				assert.Equal(t, test.want.total, got.Meta.Total)
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
				assert.Equal(t, test.want.error.Error(), errStruct.Message)
			}
		})
	}
}

func TestHandler_Read(t *testing.T) {
	type args struct {
		paramAuthorID int
	}

	type want struct {
		status   int
		usecase  *author.Schema
		response *author.GetResponse
		err      error
	}

	type test struct {
		name string
		args
		want
	}

	tests := []test{
		{
			name: "simple",
			args: args{
				paramAuthorID: 1,
			},
			want: want{
				status: http.StatusOK,
				usecase: &author.Schema{
					ID:         1,
					FirstName:  "First",
					MiddleName: "Middle",
					LastName:   "Last",
					CreatedAt:  time.Time{},
					UpdatedAt:  time.Time{},
					DeletedAt:  nil,
					Books:      nil,
				},
				response: &author.GetResponse{
					ID:         1,
					FirstName:  "First",
					MiddleName: "Middle",
					LastName:   "Last",
				},
				err: nil,
			},
		},
		{
			name: "param not supplied",
			args: args{},
			want: want{
				http.StatusBadRequest,
				&author.Schema{},
				&author.GetResponse{},
				errors.New("id is required"),
			},
		},
		{
			"simulate lower layer internal error",
			args{
				paramAuthorID: 1,
			},
			want{
				http.StatusInternalServerError,
				&author.Schema{},
				&author.GetResponse{},
				errors.New("lower layer error"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			rr := httptest.NewRequest(http.MethodGet, "/api/v1/author/{id}", nil)
			ww := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", strconv.Itoa(test.args.paramAuthorID))

			rr = rr.WithContext(context.WithValue(rr.Context(), chi.RouteCtxKey, rctx))

			router := chi.NewRouter()

			val := validator.New()

			uc := &usecase.AuthorMock{
				ReadFunc: func(ctx context.Context, authorID uint) (*author.Schema, error) {
					return test.want.usecase, test.want.err
				},
			}

			h := RegisterHTTPEndPoints(router, val, uc)
			h.Get(ww, rr)

			assert.Equal(t, test.status, ww.Code)

			if ww.Code >= 200 && ww.Code < 300 {

				var got author.GetResponse
				if err := json.NewDecoder(ww.Body).Decode(&got); err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, *test.want.response, got)
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
				assert.Equal(t, test.want.err.Error(), errStruct.Message)
			}

		})
	}
}

func TestHandler_Update(t *testing.T) {
	type args struct {
		updateRequest *author.UpdateRequest
		paramAuthorID int
	}
	type want struct {
		usecase  *author.Schema
		response *author.GetResponse
		err      error
		status   int
	}
	type test struct {
		name string
		args
		want
	}

	tests := []test{
		{
			name: "update name",
			args: args{
				updateRequest: &author.UpdateRequest{
					FirstName:  "Updated First",
					MiddleName: "Updated Middle",
					LastName:   "Updated Last",
				},
				paramAuthorID: 1,
			},
			want: want{
				usecase: &author.Schema{
					ID:         1,
					FirstName:  "Updated First",
					MiddleName: "Updated Middle",
					LastName:   "Updated Last",
					CreatedAt:  time.Time{},
					UpdatedAt:  time.Time{},
					//DeletedAt:  sql.NullTime{},
					DeletedAt: nil,
				},
				response: &author.GetResponse{
					ID:         1,
					FirstName:  "Updated First",
					MiddleName: "Updated Middle",
					LastName:   "Updated Last",
				},
				err:    nil,
				status: http.StatusOK,
			},
		},
		{
			name: "paramAuthorID not supplied",
			args: args{
				updateRequest: &author.UpdateRequest{
					FirstName:  "",
					MiddleName: "",
					LastName:   "",
				},
			},
			want: want{
				usecase:  &author.Schema{},
				response: nil,
				err:      errors.New("id is required"),
				status:   http.StatusBadRequest,
			},
		},
		{
			name: "simulate lower layer internal error",
			args: args{
				updateRequest: &author.UpdateRequest{
					FirstName: "First",
					LastName:  "Last",
				},
				paramAuthorID: 1,
			},

			want: want{
				usecase:  &author.Schema{},
				response: nil,
				err:      errors.New("lower layer error"),
				status:   http.StatusInternalServerError,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			var buf bytes.Buffer
			err := json.NewEncoder(&buf).Encode(test.args.updateRequest)
			assert.Nil(t, err)

			rr := httptest.NewRequest(http.MethodPut, "/api/v1/author/{id}", &buf)
			ww := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", strconv.Itoa(test.args.paramAuthorID))

			rr = rr.WithContext(context.WithValue(rr.Context(), chi.RouteCtxKey, rctx))

			router := chi.NewRouter()

			val := validator.New()

			uc := &usecase.AuthorMock{
				UpdateFunc: func(ctx context.Context, author *author.UpdateRequest) (*author.Schema, error) {
					return test.want.usecase, test.want.err
				}}

			h := RegisterHTTPEndPoints(router, val, uc)
			h.Update(ww, rr)

			assert.Equal(t, test.status, ww.Code)

			if ww.Code >= 200 && ww.Code < 300 {
				var got author.GetResponse
				if err = json.NewDecoder(ww.Body).Decode(&got); err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, test.want.response, &got)
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
				assert.Equal(t, test.want.err.Error(), errStruct.Message)
			}

		})
	}
}

func TestHandler_Delete(t *testing.T) {
	type args struct {
		authorID uint
	}

	type want struct {
		error
		status int
	}

	type test struct {
		name string
		args
		want
	}

	tests := []test{
		{
			name: "simple",
			args: args{
				authorID: 1,
			},
			want: want{
				error:  nil,
				status: http.StatusOK,
			},
		},
		{
			name: "paramAuthorID not provided",
			args: args{
				authorID: 0,
			},
			want: want{
				error:  errors.New("id is required"),
				status: http.StatusBadRequest,
			},
		},
		{
			name: "Simulate deleting non-existent record",
			args: args{
				authorID: 999,
			},
			want: want{
				error:  message.ErrNoRecord,
				status: http.StatusBadRequest,
			},
		},
		{
			name: "Catch-all other errors",
			args: args{
				authorID: 1,
			},
			want: want{
				error:  errors.New("all other errors"),
				status: http.StatusInternalServerError,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			rr := httptest.NewRequest(http.MethodDelete, "/api/v1/author/{id}", nil)
			ww := httptest.NewRecorder()

			rctx := chi.NewRouteContext()
			rctx.URLParams.Add("id", strconv.Itoa(int(test.args.authorID)))

			rr = rr.WithContext(context.WithValue(rr.Context(), chi.RouteCtxKey, rctx))

			router := chi.NewRouter()
			val := validator.New()

			uc := &usecase.AuthorMock{
				DeleteFunc: func(ctx context.Context, authorID uint) error {
					return test.want.error
				},
			}

			h := RegisterHTTPEndPoints(router, val, uc)
			h.Delete(ww, rr)

			assert.Equal(t, test.status, ww.Code)
		})
	}
}
