package usecase

import (
	"context"

	"github.com/gmhafiz/go8/internal/domain/book"
	"github.com/gmhafiz/go8/internal/domain/book/repository"
)

//go:generate mirip -rm -pkg usecase -out usecase_mock.go . Book
type Book interface {
	Create(ctx context.Context, book *book.CreateRequest) (*book.DB, error)
	List(ctx context.Context, f *book.Filter) ([]*book.DB, error)
	Read(ctx context.Context, bookID int) (*book.DB, error)
	Update(ctx context.Context, book *book.UpdateRequest) (*book.DB, error)
	Delete(ctx context.Context, bookID int) error
	Search(ctx context.Context, req *book.Filter) ([]*book.DB, error)
}

type BookUseCase struct {
	bookRepo repository.Book
}

func New(bookRepo repository.Book) *BookUseCase {
	return &BookUseCase{
		bookRepo: bookRepo,
	}
}

func (u *BookUseCase) Create(ctx context.Context, book *book.CreateRequest) (*book.DB, error) {
	bookID, err := u.bookRepo.Create(ctx, book)
	if err != nil {
		return nil, err
	}
	bookFound, err := u.bookRepo.Read(context.Background(), bookID)
	if err != nil {
		return nil, err
	}
	return bookFound, err
}

func (u *BookUseCase) List(ctx context.Context, f *book.Filter) ([]*book.DB, error) {
	return u.bookRepo.List(ctx, f)
}

func (u *BookUseCase) Read(ctx context.Context, bookID int) (*book.DB, error) {
	return u.bookRepo.Read(ctx, bookID)
}

func (u *BookUseCase) Update(ctx context.Context, book *book.UpdateRequest) (*book.DB, error) {
	err := u.bookRepo.Update(ctx, book)
	if err != nil {
		return nil, err
	}
	return u.bookRepo.Read(ctx, book.ID)
}

func (u *BookUseCase) Delete(ctx context.Context, bookID int) error {
	return u.bookRepo.Delete(ctx, bookID)
}

func (u *BookUseCase) Search(ctx context.Context, req *book.Filter) ([]*book.DB, error) {
	return u.bookRepo.Search(ctx, req)
}
