package usecase

import (
	"context"
	"testing"

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

func TestBookUseCase_Create(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	repo := mock.NewMockRepository(ctrl)
	uc := NewBookUseCase(repo)

	request := book.Request{
		Title:         "title",
		PublishedDate: "2006-01-02 15:04:05 +0000 UTC",
		ImageURL:      "https://example.com/image.png",
		Description:   "",
	}

	ctx := context.Background()

	expected := &models.Book{
		BookID:        0,
		Title:         request.Title,
		PublishedDate: now.MustParse(request.PublishedDate),
		ImageURL: null.String{
			String: request.ImageURL,
			Valid:  true,
		},
		Description: request.Description,
	}
	var err error
	var bookID int64
	repo.EXPECT().Create(ctx, gomock.Eq(expected)).Return(bookID, err).AnyTimes()
	repo.EXPECT().Find(ctx, gomock.Any()).Return(expected, err).AnyTimes()

	bookGot, err := uc.Create(ctx, request)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotEqual(t, bookGot.BookID, 0)
	assert.Equal(t, bookGot.Title, request.Title)
	assert.Equal(t, bookGot.PublishedDate.String(), request.PublishedDate)
	assert.Equal(t, bookGot.Description, request.Description)
	assert.Equal(t, bookGot.ImageURL.String, request.ImageURL)
}
