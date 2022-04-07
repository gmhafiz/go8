package handler

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"github.com/stretchr/testify/assert"

	"github.com/gmhafiz/go8/ent/gen"
	"github.com/gmhafiz/go8/internal/domain/author"
	"github.com/gmhafiz/go8/internal/domain/author/usecase"
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
		*gen.Author
		error
		Errs
	}

	type test struct {
		name string
		args
		want
		status int
	}

	tests := []test{
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
				Author: &gen.Author{
					ID:         1,
					FirstName:  "First",
					MiddleName: "Middle",
					LastName:   "Last",
					CreatedAt:  time.Time{},
					UpdatedAt:  time.Time{},
					DeletedAt:  nil,
					Edges: gen.AuthorEdges{
						Books: nil,
					},
				},
				error: nil,
			},
			status: http.StatusCreated,
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
				Author: &gen.Author{},
				Errs: Errs{
					Message: []string{"CreateRequest.FirstName is required"},
				},
			},
			status: http.StatusBadRequest,
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
				Author: &gen.Author{},
				error:  ErrTransactionFailed,
			},
			status: http.StatusInternalServerError,
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
				Author: &gen.Author{},
				error:  sql.ErrNoRows,
			},
			status: http.StatusBadRequest,
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
				Author: &gen.Author{},
				error:  errors.New("other error"),
			},
			status: http.StatusInternalServerError,
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
				CreateFunc: func(ctx context.Context, a *author.CreateRequest) (*gen.Author, error) {
					return test.want.Author, test.want.error
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
				var got gen.Author
				if err = json.NewDecoder(ww.Body).Decode(&got); err != nil {
					t.Fatal(err)
				}

				assert.Equal(t, ww.Code, test.status)
				assert.Equal(t, &got, test.want.Author)
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
			authors []*gen.Author
			total   int
			error
		}
		status int
		size   int
		total  int
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
				status: http.StatusOK,
				size:   0,
				total:  0,
				usecase: struct {
					authors []*gen.Author
					total   int
					error
				}{
					authors: make([]*gen.Author, 0),
					total:   0,
					error:   nil,
				},
			},
		},
		{
			name: "returns one record",
			args: args{
				uri: "/api/v1/author",
			},
			want: want{
				status: http.StatusOK,
				size:   1,
				total:  1,
				usecase: struct {
					authors []*gen.Author
					total   int
					error
				}{
					authors: []*gen.Author{
						{
							ID:         1,
							FirstName:  "First",
							MiddleName: "Middle",
							LastName:   "Last",
							CreatedAt:  time.Time{},
							UpdatedAt:  time.Time{},
							DeletedAt:  nil,
							Edges: gen.AuthorEdges{
								Books: nil,
							},
						},
					},
					total: 1,
					error: nil,
				},
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
					authors []*gen.Author
					total   int
					error
				}{
					authors: []*gen.Author{},
					total:   0,
					error:   errors.New("some use case internal error"),
				},
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
				ListFunc: func(ctx context.Context, f *author.Filter) ([]*gen.Author, int, error) {
					return test.want.usecase.authors, test.want.usecase.total, test.want.usecase.error
				},
			}

			h := RegisterHTTPEndPoints(router, val, uc)
			h.List(ww, rr)

			var got respond.Standard
			if err := json.NewDecoder(ww.Body).Decode(&got); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, test.want.status, ww.Code)
			assert.Equal(t, test.want.size, got.Meta.Size)
			assert.Equal(t, test.want.total, got.Meta.Total)
		})
	}
}

func TestHandler_Read(t *testing.T) {
	type args struct {
		paramAuthorID int
	}

	type want struct {
		*gen.Author
		error
	}

	type test struct {
		name string
		args
		want
		status int
	}

	tests := []test{
		{
			name: "simple",
			args: args{
				paramAuthorID: 1,
			},
			want: want{
				Author: &gen.Author{
					ID:         1,
					FirstName:  "First",
					MiddleName: "Middle",
					LastName:   "Last",
					CreatedAt:  time.Time{},
					UpdatedAt:  time.Time{},
					DeletedAt:  nil,
					Edges: gen.AuthorEdges{
						Books: nil,
					},
				},
				error: nil,
			},
			status: http.StatusOK,
		},
		{
			"param not supplied",
			args{},
			want{
				&gen.Author{},
				errors.New("id is required"),
			},
			http.StatusBadRequest,
		},
		{
			"simulate lower layer internal error",
			args{paramAuthorID: 1},
			want{
				&gen.Author{},
				errors.New("lower layer error"),
			},
			http.StatusInternalServerError,
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
				ReadFunc: func(ctx context.Context, authorID uint) (*gen.Author, error) {
					return test.want.Author, test.want.error
				},
			}

			h := RegisterHTTPEndPoints(router, val, uc)
			h.Get(ww, rr)

			var got gen.Author
			if err := json.NewDecoder(ww.Body).Decode(&got); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, test.status, ww.Code)
			assert.Equal(t, *test.want.Author, got)
		})
	}
}

func TestHandler_Update(t *testing.T) {
	type args struct {
		updateRequest *author.Update
		paramAuthorID int
	}
	type want struct {
		*gen.Author
		error
	}
	type test struct {
		name string
		args
		want
		status int
	}

	tests := []test{
		{
			name: "update name",
			args: args{
				updateRequest: &author.Update{
					FirstName:  "Updated First",
					MiddleName: "Updated Middle",
					LastName:   "Updated Last",
				},
				paramAuthorID: 1,
			},
			want: want{
				Author: &gen.Author{
					ID:         1,
					FirstName:  "Updated First",
					MiddleName: "Updated Middle",
					LastName:   "Updated Last",
					CreatedAt:  time.Time{},
					UpdatedAt:  time.Time{},
					DeletedAt:  nil,
				},
				error: nil,
			},
			status: http.StatusOK,
		},
		{
			name: "paramAuthorID not supplied",
			args: args{
				updateRequest: &author.Update{
					FirstName:  "",
					MiddleName: "",
					LastName:   "",
				},
			},
			want: want{
				Author: &gen.Author{},
				error:  errors.New("id is required"),
			},
			status: http.StatusBadRequest,
		},
		{
			name: "simulate lower layer internal error",
			args: args{
				updateRequest: &author.Update{
					FirstName: "First",
					LastName:  "Last",
				},
				paramAuthorID: 1,
			},

			want: want{
				Author: &gen.Author{},
				error:  errors.New("lower layer error"),
			},

			status: http.StatusInternalServerError,
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
				UpdateFunc: func(ctx context.Context, author *author.Update) (*gen.Author, error) {
					return test.want.Author, test.want.error
				}}

			h := RegisterHTTPEndPoints(router, val, uc)
			h.Update(ww, rr)

			var got gen.Author
			if err = json.NewDecoder(ww.Body).Decode(&got); err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, test.status, ww.Code)
			assert.Equal(t, test.want.Author, &got)
		})
	}
}

func TestHandler_Delete(t *testing.T) {
	type args struct {
		authorID uint
	}

	type want struct {
		error
	}

	type test struct {
		name string
		args
		want
		status int
	}

	tests := []test{
		{
			name: "simple",
			args: args{
				authorID: 1,
			},
			want: want{
				error: nil,
			},
			status: http.StatusOK,
		},
		{
			name: "paramAuthorID not provided",
			args: args{
				authorID: 0,
			},
			want: want{
				error: errors.New("id is required"),
			},
			status: http.StatusBadRequest,
		},
		{
			name: "Simulate deleting non-existent record",
			args: args{
				authorID: 999,
			},
			want: want{
				error: respond.ErrNoRecord,
			},
			status: http.StatusBadRequest,
		},
		{
			name: "Catch-all other errors",
			args: args{
				authorID: 1,
			},
			want: want{
				error: errors.New("all other errors"),
			},
			status: http.StatusInternalServerError,
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
