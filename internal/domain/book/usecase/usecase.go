package usecase

import (
	"context"

	"github.com/gmhafiz/go8/internal/domain/book"
	"github.com/gmhafiz/go8/internal/models"
)

type BookUseCase struct {
	bookRepo book.Repository
}

func NewBookUseCase(bookRepo book.Repository) *BookUseCase {
	return &BookUseCase{
		bookRepo: bookRepo,
	}
}

func (u *BookUseCase) Create(ctx context.Context, r book.Request) (*models.Book, error) {
	bk := book.ToBook(&r)
	bookID, err := u.bookRepo.Create(ctx, bk)
	if err != nil {
		return nil, err
	}
	bookFound, err := u.bookRepo.Find(context.Background(), bookID)
	if err != nil {
		return nil, err
	}
	return bookFound, err
}

func (u *BookUseCase) All(ctx context.Context) ([]*models.Book, error) {
	return u.bookRepo.All(ctx)
}

func (u *BookUseCase) Find(ctx context.Context, bookID int64) (*models.Book, error) {
	return u.bookRepo.Find(ctx, bookID)
}

func (u *BookUseCase) Update(ctx context.Context, book *models.Book) (*models.Book, error) {
	return u.bookRepo.Update(ctx, book)
}

func (u *BookUseCase) Delete(ctx context.Context, bookID int64) error {
	return u.bookRepo.Delete(ctx, bookID)
}

func (u *BookUseCase) Search(ctx context.Context, req *book.Request) ([]*models.Book, error) {
	return u.bookRepo.Search(ctx, req)
}
