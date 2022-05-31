package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/gmhafiz/go8/config"
	"github.com/gmhafiz/go8/ent/gen"
	"github.com/gmhafiz/go8/internal/domain/author"
	"github.com/gmhafiz/go8/internal/domain/author/repository"
	"github.com/gmhafiz/go8/internal/utility/filter"
)

var c config.Cache

func TestMain(m *testing.M) {
	c = config.Cache{
		Enable: false,
	}
}

func TestAuthorUseCase_Create(t *testing.T) {
	type args struct {
		*author.CreateRequest
	}

	type want struct {
		*gen.Author
		error
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
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			repoAuthor := &repository.AuthorMock{
				CreateFunc: func(ctx context.Context, r *author.CreateRequest) (*gen.Author, error) {
					return test.want.Author, test.want.error
				},
			}

			uc := New(c, repoAuthor, nil, nil, nil)

			got, err := uc.Create(context.Background(), test.args.CreateRequest)
			assert.Equal(t, test.want.error, err)
			assert.Equal(t, test.want.Author, got)
		})
	}
}

func TestAuthorUseCase_List(t *testing.T) {
	type args struct {
		context.Context
		filter *author.Filter
	}
	type want struct {
		authors []*gen.Author
		total   int
		error
	}
	type test struct {
		name string
		args
		want
	}

	twoAuthors := []*gen.Author{
		{
			ID:         1,
			FirstName:  "1 First",
			MiddleName: "",
			LastName:   "2 Last",
			CreatedAt:  time.Time{},
			UpdatedAt:  time.Time{},
			DeletedAt:  &time.Time{},
			Edges: gen.AuthorEdges{
				Books: nil,
			},
		},
		{
			ID:         2,
			FirstName:  "2 First",
			MiddleName: "",
			LastName:   "2 Last",
			CreatedAt:  time.Time{},
			UpdatedAt:  time.Time{},
			DeletedAt:  &time.Time{},
			Edges: gen.AuthorEdges{
				Books: nil,
			},
		},
	}

	searched := []*gen.Author{
		{
			ID:         2,
			FirstName:  "2 First",
			MiddleName: "",
			LastName:   "2 Last",
			CreatedAt:  time.Time{},
			UpdatedAt:  time.Time{},
			DeletedAt:  &time.Time{},
			Edges: gen.AuthorEdges{
				Books: nil,
			},
		},
	}

	tests := []test{
		{
			name: "simple",
			args: args{
				Context: context.Background(),
				filter: &author.Filter{
					Base: filter.Filter{
						Page:          0,
						Offset:        0,
						Limit:         10,
						DisablePaging: false,
						Sort:          nil,
						Search:        false,
					},
					FirstName:  "",
					MiddleName: "",
					LastName:   "",
				},
			},
			want: want{
				authors: twoAuthors,
				total:   2,
				error:   nil,
			},
		},
		{
			name: "search",
			args: args{
				Context: context.Background(),
				filter: &author.Filter{
					Base: filter.Filter{
						Page:          0,
						Offset:        0,
						Limit:         0,
						DisablePaging: false,
						Sort:          nil,
						Search:        true,
					},
					FirstName:  "2 First",
					MiddleName: "",
					LastName:   "",
				},
			},
			want: want{
				authors: searched,
				total:   1,
				error:   nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			repoAuthor := &repository.AuthorMock{
				ListFunc: func(ctx context.Context, f *author.Filter) ([]*gen.Author, int, error) {
					return test.want.authors, test.want.total, test.want.error
				},
			}

			cacheMock := &repository.AuthorRedisServiceMock{
				ListFunc: func(ctx context.Context, f *author.Filter) ([]*gen.Author, int, error) {
					return test.want.authors, test.want.total, test.want.error
				},
			}

			searchMock := &repository.SearcherMock{
				SearchFunc: func(ctx context.Context, f *author.Filter) ([]*gen.Author, int, error) {
					return test.want.authors, test.want.total, test.want.error
				},
			}

			uc := New(c, repoAuthor, searchMock, nil, cacheMock)

			got, total, err := uc.List(test.args.Context, test.args.filter)
			assert.Equal(t, test.want.error, err)
			assert.Equal(t, test.want.total, total)
			assert.Equal(t, test.want.authors, got)
		})
	}
}

func TestAuthorUseCase_Read(t *testing.T) {
	type args struct {
		ID uint
	}

	type want struct {
		*gen.Author
		error
	}

	type test struct {
		name string
		args
		want
	}

	tests := []test{
		{
			name: "one",
			args: args{
				ID: 1,
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
		},
		{
			name: "zero ID",
			args: args{
				ID: 0,
			},
			want: want{
				Author: nil,
				error:  errors.New("ID cannot be 0"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			repoAuthor := &repository.AuthorMock{
				ReadFunc: func(ctx context.Context, id uint) (*gen.Author, error) {
					return test.want.Author, test.want.error
				},
			}

			uc := New(c, repoAuthor, nil, nil, nil)

			got, err := uc.Read(context.Background(), test.args.ID)
			assert.Equal(t, test.want.error, err)

			assert.Equal(t, test.want.Author, got)
		})
	}
}

func TestAuthorUseCase_Update(t *testing.T) {
	type args struct {
		context.Context
		*author.Update
	}
	type want struct {
		repo struct {
			*gen.Author
			error
		}
		error
	}

	type test struct {
		name string
		args
		want
	}

	createdTime := time.Now()

	tests := []test{
		{
			name: "simple",
			args: args{
				Context: context.Background(),
				Update: &author.Update{
					ID:         1,
					FirstName:  "Updated First",
					MiddleName: "Updated Middle",
					LastName:   "Updated Last",
				},
			},
			want: want{
				repo: struct {
					*gen.Author
					error
				}{
					&gen.Author{
						ID:         1,
						FirstName:  "Updated First",
						MiddleName: "Updated Middle",
						LastName:   "Updated Last",
						CreatedAt:  createdTime,
						UpdatedAt:  time.Now(),
						DeletedAt:  nil,
						Edges: gen.AuthorEdges{
							Books: nil,
						},
					},
					nil,
				},
				error: nil,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			repoAuthor := &AuthorMock{
				UpdateFunc: func(ctx context.Context, authorMiripParam *author.Update) (*gen.Author, error) {
					return test.want.repo.Author, test.want.repo.error
				},
			}

			cacheMock := &repository.AuthorRedisServiceMock{
				UpdateFunc: func(ctx context.Context, toAuthor *author.Update) (*gen.Author, error) {
					return test.want.repo.Author, test.want.repo.error
				},
			}

			uc := New(c, repoAuthor, nil, nil, cacheMock)

			update, err := uc.Update(test.args.Context, test.args.Update)
			assert.Equal(t, test.want.error, err)

			assert.Equal(t, test.want.repo.ID, update.ID)
			assert.Equal(t, test.want.repo.FirstName, update.FirstName)
			assert.Equal(t, test.want.repo.MiddleName, update.MiddleName)
			assert.Equal(t, test.want.repo.LastName, update.LastName)
			assert.True(t, createdTime.Before(test.want.repo.CreatedAt) || createdTime.Equal(test.want.repo.CreatedAt))
			assert.True(t, createdTime.Before(test.want.repo.UpdatedAt) || createdTime.Equal(test.want.repo.UpdatedAt))
			assert.Nil(t, test.want.repo.DeletedAt)
		})
	}
}

func TestAuthorUseCase_Delete(t *testing.T) {
	type args struct {
		context.Context
		ID uint
	}
	type want struct {
		error
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
				Context: context.Background(),
				ID:      1,
			},
			want: want{
				error: nil,
			},
		},
		{
			name: "zero ID",
			args: args{
				Context: context.Background(),
				ID:      0,
			},
			want: want{
				error: errors.New("ID cannot be 0 or less"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			repoAuthor := &AuthorMock{
				DeleteFunc: func(ctx context.Context, authorID uint) error {
					return test.want.error
				},
			}
			cacheMock := &repository.AuthorRedisServiceMock{
				DeleteFunc: func(ctx context.Context, id uint) error {
					return test.want.error
				},
			}

			uc := New(c, repoAuthor, nil, nil, cacheMock)

			err := uc.Delete(test.args.Context, test.args.ID)
			assert.Equal(t, test.want.error, err)
		})
	}
}
