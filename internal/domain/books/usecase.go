package books

import (
	"context"

	"go8ddd/internal/model"
)

type BookUseCase interface {
	All(ctx context.Context) (model.BookSlice, error)
	Create(ctx context.Context, book *model.Book) (*model.Book, error)
	Get(ctx context.Context, bookID int64) (*model.Book, error)
	Update(ctx context.Context, book *model.Book) (*model.Book, error)
	Delete(ctx context.Context, bookID int64, hardDelete bool) error
}

type bookUseCase struct {
	bookRepo BookRepository
}

func (u *bookUseCase) All(ctx context.Context) (model.BookSlice, error) {
	bookList, err := u.bookRepo.All(ctx)
	if err != nil {
		return nil, err
	}

	return bookList, nil
}

func (u *bookUseCase) Create(ctx context.Context, book *model.Book) (*model.Book, error) {
	book, err := u.bookRepo.Create(ctx, book)
	if err != nil {
		return nil, err
	}

	return book, nil
}

func (u *bookUseCase) Get(ctx context.Context, bookID int64) (*model.Book, error) {
	book, err := u.bookRepo.Get(ctx, bookID)
	if err != nil {
		return nil, err
	}

	return book, nil
}

func (u *bookUseCase) Update(ctx context.Context, book *model.Book) (*model.Book, error) {
	return u.bookRepo.Update(ctx, book)
}

func (u *bookUseCase) Delete(ctx context.Context, bookID int64, hardDelete bool) error {
	if hardDelete {
		return u.bookRepo.HardDelete(ctx, bookID)
	}
	return u.bookRepo.Delete(ctx, bookID)
}

func NewUseCase(repo BookRepository) BookUseCase {
	return &bookUseCase{
		bookRepo: repo,
	}
}
