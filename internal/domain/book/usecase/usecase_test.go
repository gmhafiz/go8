package usecase

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/jinzhu/now"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"

	"github.com/gmhafiz/go8/internal/domain/book"
	"github.com/gmhafiz/go8/internal/domain/book/mock"
	"github.com/gmhafiz/go8/internal/models"
)

//go:generate mockgen -package mock -source ../usecase.go -destination=../mock/mock_usecase.go

func newUseCase(t *testing.T) (*BookUseCase, *mock.MockRepository) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockRepository(ctrl)
	return New(repo), repo
}

func TestBookUseCase_Create(t *testing.T) {
	uc, repo := newUseCase(t)

	request := &models.Book{
		Title:         "title",
		PublishedDate: now.MustParse("2006-01-02 15:04:05 +0000 UTC"),
		ImageURL: null.String{
			String: "https://example.com/image.png",
			Valid:  true,
		},
		Description: "",
	}

	ctx := context.Background()

	expected := &models.Book{
		BookID:        0,
		Title:         request.Title,
		PublishedDate: request.PublishedDate,
		ImageURL: null.String{
			String: request.ImageURL.String,
			Valid:  true,
		},
		Description: request.Description,
	}
	var err error
	var bookID int64
	repo.EXPECT().Create(ctx, gomock.Eq(expected)).Return(bookID, err).AnyTimes()
	repo.EXPECT().Read(ctx, gomock.Any()).Return(expected, err).AnyTimes()

	bookGot, err := uc.Create(ctx, request)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEqual(t, bookGot.BookID, 0)
	assert.Equal(t, bookGot.Title, request.Title)
	assert.Equal(t, bookGot.PublishedDate.String(), request.PublishedDate.String())
	assert.Equal(t, bookGot.Description, request.Description)
	assert.Equal(t, bookGot.ImageURL.String, request.ImageURL.String)
}

func TestBookUseCase_List(t *testing.T) {
	uc, repo := newUseCase(t)

	ctx := context.Background()
	var err error
	var want []*models.Book
	filter := &book.Filter{}

	repo.EXPECT().List(ctx, filter).Return(want, err).AnyTimes()

	books, err := uc.List(ctx, filter)

	assert.NoError(t, err)
	assert.Nil(t, books)
}

func TestBookUseCase_Read(t *testing.T) {
	uc, repo := newUseCase(t)

	ctx := context.Background()
	var err error
	var id int64
	var want *models.Book

	repo.EXPECT().Read(ctx, id).Return(want, err).AnyTimes()

	_, err = uc.Read(ctx, id)

	assert.NoError(t, err)
}

func TestBookUseCase_Update(t *testing.T) {
	uc, repo := newUseCase(t)
	ctx := context.Background()
	var err error

	request := &models.Book{
		BookID:        1,
		Title:         "updated title",
		PublishedDate: time.Time{},
		ImageURL:      null.String{},
		Description:   "",
	}

	repo.EXPECT().Update(ctx, request).Return(err).AnyTimes()
	repo.EXPECT().Read(ctx, gomock.Any()).Return(request, err).AnyTimes()

	got, err := uc.Update(ctx, request)

	assert.NoError(t, err)
	assert.Equal(t, request.BookID, got.BookID)
	assert.Equal(t, request.Title, got.Title)
	assert.Equal(t, request.Description, got.Description)
}

func TestBookUseCase_Delete(t *testing.T) {
	uc, repo := newUseCase(t)
	ctx := context.Background()
	var id int64

	repo.EXPECT().Delete(ctx, gomock.Any()).Return(nil).AnyTimes()

	err := uc.Delete(ctx, id)

	assert.NoError(t, err)
}

func TestBookUseCase_Search(t *testing.T) {
	uc, repo := newUseCase(t)
	ctx := context.Background()
	var err error
	var want []*models.Book
	filter := &book.Filter{}

	repo.EXPECT().Search(ctx, filter).Return(want, err).AnyTimes()

	got, err := uc.Search(ctx, filter)

	assert.NoError(t, err)
	assert.Len(t, got, 0)
}
