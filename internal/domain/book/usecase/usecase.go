package usecase

import (
	"context"

	"github.com/jinzhu/now"
	"github.com/volatiletech/null/v8"

	"github.com/gmhafiz/go8/internal/domain/book"
	"github.com/gmhafiz/go8/internal/model"
)

type BookUseCase struct {
	bookRepo book.Repository
}

func NewBookUseCase(bookRepo book.Repository) *BookUseCase {
	return &BookUseCase{
		bookRepo: bookRepo,
	}
}

func (b *BookUseCase) Create(ctx context.Context, title, description, imageURL, publishedDate string) (*model.Book, error) {
	bk := &model.Book{
		Title:         title,
		PublishedDate: now.MustParse(publishedDate),
		ImageURL: null.String{
			String: imageURL,
			Valid:  true,
		},
		Description: null.String{
			String: description,
			Valid:  true,
		},
	}

	bookID, err := b.bookRepo.Create(ctx, bk)
	if err != nil {
		return nil, err
	}
	bookFound, err := b.bookRepo.Find(context.Background(), bookID)
	if err != nil {
		return nil, err
	}
	return bookFound, err
}

func (b *BookUseCase) All(ctx context.Context) ([]*model.Book, error) {
	return b.bookRepo.All(ctx)
}

func (b *BookUseCase) Find(ctx context.Context, bookID int64) (*model.Book, error) {
	return b.bookRepo.Find(ctx, bookID)
}
