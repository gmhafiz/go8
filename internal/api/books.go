package api

import (
	"context"

	"eight/internal/models"
)

func (a API) GetAllBooks(ctx context.Context) (models.BookSlice, error) {
	return a.books.AllBooks(ctx)
}

func (a API) CreateBook(ctx context.Context, book *models.Book) (*models.Book, error) {
	return a.books.CreateBook(ctx, book)
}

func (a API) GetBook(ctx context.Context, bookID int64) (*models.Book, error) {
	return a.books.GetBook(ctx, bookID)
}

func (a API) Delete(ctx context.Context, bookID int64) error {
	return a.books.Delete(ctx, bookID)
}
