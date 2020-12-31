package usecase

import (
	"context"
	"time"

	"github.com/volatiletech/null/v8"

	"github.com/gmhafiz/go8/internal/domain/book"
	"github.com/gmhafiz/go8/internal/model"
	"github.com/gmhafiz/go8/internal/resource"
)

type BookUseCase struct {
	bookRepo book.Repository
}

func NewBookUseCase(bookRepo book.Repository) *BookUseCase {
	return &BookUseCase{
		bookRepo: bookRepo,
	}
}

func (b *BookUseCase) Create(ctx context.Context, title, description string) (*model.Book, error) {
	bk := &model.Book{
		Title: title,
		Description: null.String{
			String: description,
			Valid:  true,
		},
		PublishedDate: time.Now(),
	}

	bk, err := b.bookRepo.Create(ctx, bk)
	return bk, err
}

func (b *BookUseCase) All(ctx context.Context) ([]resource.BookDB, error) {
	return b.bookRepo.All(ctx)
}
