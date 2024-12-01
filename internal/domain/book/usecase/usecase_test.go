package usecase

import (
	"context"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"

	"github.com/gmhafiz/go8/internal/domain/book"
	"github.com/gmhafiz/go8/internal/domain/book/repository"
	"github.com/gmhafiz/go8/internal/utility/filter"
)

func TestBookUseCase_Create(t *testing.T) {
	type args struct {
		ctx context.Context
		req *book.CreateRequest
	}

	type want struct {
		book *book.Schema
		err  error
	}

	type test struct {
		name string
		args
		want
		*repository.BookMock
	}

	timeParsed, err := time.Parse(time.RFC3339, "2020-02-02T00:00:00Z")
	assert.Nil(t, err)

	tests := []test{
		{
			name: "simple",
			args: args{
				ctx: context.Background(),
				req: &book.CreateRequest{
					Title:         "title",
					PublishedDate: "2020-02-02T00:00:00Z",
					ImageURL:      "https://example.com/image.png",
					Description:   "description",
				},
			},
			want: want{
				book: &book.Schema{
					ID:            1,
					Title:         "title",
					PublishedDate: timeParsed,
					ImageURL:      "https://example.com/image.png",
					Description:   "description",
				},
				err: nil,
			},
			BookMock: &repository.BookMock{
				CreateFunc: func(ctx context.Context, bookMiripParam *book.CreateRequest) (uint64, error) {
					return 1, nil
				},
				ReadFunc: func(ctx context.Context, bookID uint64) (*book.Schema, error) {
					return &book.Schema{
						ID:            1,
						Title:         "title",
						PublishedDate: timeParsed,
						ImageURL:      "https://example.com/image.png",
						Description:   "description",
					}, nil
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			uc := New(test.BookMock)

			created, err := uc.Create(test.args.ctx, test.args.req)
			assert.Equal(t, test.want.err, err)

			assert.Equal(t, test.want.book.ID, created.ID)
			assert.Equal(t, test.want.book.Title, created.Title)
			assert.Equal(t, test.want.book.PublishedDate, created.PublishedDate)
			assert.Equal(t, test.want.book.ImageURL, created.ImageURL)
			assert.Equal(t, test.want.book.Description, created.Description)
		})
	}
}

func TestBookUseCase_List(t *testing.T) {
	type fields struct {
		bookRepo repository.BookMock
	}
	type args struct {
		ctx context.Context
		f   *book.Filter
	}

	timeParsed, err := time.Parse(time.RFC3339, "2020-02-02T00:00:00Z")
	assert.Nil(t, err)

	oneBook := []*book.Schema{
		{
			ID:            1,
			Title:         "title 1",
			PublishedDate: timeParsed,
			ImageURL:      "https://example.com/image1.png",
			Description:   "description 1",
		},
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*book.Schema
		wantErr error
	}{
		{
			name: "simple",
			fields: fields{
				bookRepo: repository.BookMock{
					ListFunc: func(ctx context.Context, f *book.Filter) ([]*book.Schema, error) {
						return oneBook, nil
					},
				},
			},
			args: args{
				ctx: context.Background(),
				f: &book.Filter{
					Base: filter.Filter{
						Page:          1,
						Offset:        0,
						Limit:         0,
						DisablePaging: false,
						Sort:          nil,
						Search:        false,
					},
					Title:         "",
					Description:   "",
					PublishedDate: "",
				},
			},
			want:    oneBook,
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &BookUseCase{
				bookRepo: &tt.fields.bookRepo,
			}
			got, err := u.List(tt.args.ctx, tt.args.f)
			assert.Equal(t, tt.wantErr, err)
			assert.Equalf(t, tt.want, got, "List(%v, %v)", tt.args.ctx, tt.args.f)
		})
	}
}

func TestBookUseCase_Read(t *testing.T) {
	type fields struct {
		bookRepo repository.Book
	}
	type args struct {
		ctx    context.Context
		bookID uint64
	}

	timeParsed, err := time.Parse(time.RFC3339, "2020-02-02T00:00:00Z")
	assert.Nil(t, err)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *book.Schema
		wantErr error
	}{
		{
			name: "simple",
			fields: fields{
				bookRepo: &repository.BookMock{
					ReadFunc: func(ctx context.Context, bookID uint64) (*book.Schema, error) {
						return &book.Schema{
							ID:            1,
							Title:         "title",
							PublishedDate: timeParsed,
							ImageURL:      "https://example.com/image.png",
							Description:   "description",
						}, nil
					}},
			},
			args: args{
				ctx:    context.Background(),
				bookID: 1,
			},
			want: &book.Schema{
				ID:            1,
				Title:         "title",
				PublishedDate: timeParsed,
				ImageURL:      "https://example.com/image.png",
				Description:   "description",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &BookUseCase{
				bookRepo: tt.fields.bookRepo,
			}
			got, err := u.Read(tt.args.ctx, tt.args.bookID)
			assert.Equal(t, err, tt.wantErr)
			assert.Equalf(t, tt.want, got, "Read(%v, %v)", tt.args.ctx, tt.args.bookID)
		})
	}
}

func TestBookUseCase_Update(t *testing.T) {
	type fields struct {
		bookRepo repository.Book
	}
	type args struct {
		ctx  context.Context
		book *book.UpdateRequest
	}

	timeParsed, err := time.Parse(time.RFC3339, "2020-02-02T00:00:00Z")
	assert.Nil(t, err)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *book.Schema
		wantErr error
	}{
		{
			name: "simple",
			fields: fields{
				bookRepo: &repository.BookMock{
					UpdateFunc: func(ctx context.Context, book *book.UpdateRequest) error {
						return nil
					},
					ReadFunc: func(ctx context.Context, bookID uint64) (*book.Schema, error) {
						return &book.Schema{
							ID:            1,
							Title:         "title",
							PublishedDate: timeParsed,
							ImageURL:      "https://example.com/image1.png",
							Description:   "description",
						}, nil
					},
				},
			},
			args: args{
				ctx: context.Background(),
				book: &book.UpdateRequest{
					Title:         "title",
					PublishedDate: "2020-02-02T00:00:00Z",
					ImageURL:      "https://example.com/image1.png",
					Description:   "description",
				},
			},
			want: &book.Schema{
				ID:            1,
				Title:         "title",
				PublishedDate: timeParsed,
				ImageURL:      "https://example.com/image1.png",
				Description:   "description",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &BookUseCase{
				bookRepo: tt.fields.bookRepo,
			}
			got, err := u.Update(tt.args.ctx, tt.args.book)
			assert.Equal(t, tt.wantErr, err)
			assert.Equalf(t, tt.want, got, "Update(%v, %v)", tt.args.ctx, tt.args.book)
		})
	}
}

func TestBookUseCase_Delete(t *testing.T) {
	type fields struct {
		bookRepo repository.Book
	}
	type args struct {
		ctx    context.Context
		bookID uint64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{
			name: "simple",
			fields: fields{
				bookRepo: &repository.BookMock{
					DeleteFunc: func(ctx context.Context, bookID uint64) error {
						return nil
					},
				},
			},
			args: args{
				ctx:    context.Background(),
				bookID: 1,
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &BookUseCase{
				bookRepo: tt.fields.bookRepo,
			}
			err := u.Delete(tt.args.ctx, tt.args.bookID)
			assert.Equal(t, tt.wantErr, err)
		})
	}
}

func TestBookUseCase_Search(t *testing.T) {
	type fields struct {
		bookRepo repository.Book
	}
	type args struct {
		ctx context.Context
		req *book.Filter
	}

	timeParsed, err := time.Parse(time.RFC3339, "2020-02-02T00:00:00Z")
	assert.Nil(t, err)

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*book.Schema
		wantErr error
	}{
		{
			name: "simple",
			fields: fields{
				bookRepo: &repository.BookMock{
					SearchFunc: func(ctx context.Context, req *book.Filter) ([]*book.Schema, error) {
						return []*book.Schema{
							{
								ID:            1,
								Title:         "searched 1",
								PublishedDate: timeParsed,
								ImageURL:      "https://example.com/image1.png",
								Description:   "description",
							},
						}, nil
					},
				},
			},
			args: args{
				ctx: nil,
				req: &book.Filter{
					Base: filter.Filter{
						Page:          1,
						Offset:        0,
						Limit:         0,
						DisablePaging: false,
						Sort:          nil,
						Search:        true,
					},
					Title:         "searched",
					Description:   "",
					PublishedDate: "",
				},
			},
			want: []*book.Schema{
				{
					ID:            1,
					Title:         "searched 1",
					PublishedDate: timeParsed,
					ImageURL:      "https://example.com/image1.png",
					Description:   "description",
				},
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			u := &BookUseCase{
				bookRepo: tt.fields.bookRepo,
			}
			got, err := u.Search(tt.args.ctx, tt.args.req)
			assert.Equal(t, tt.wantErr, err)
			assert.Equalf(t, tt.want, got, "Search(%v, %v)", tt.args.ctx, tt.args.req)
		})
	}
}
