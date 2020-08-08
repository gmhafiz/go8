package books

import (
	"context"
	"database/sql"

	"github.com/jackc/pgx/v4/pgxpool"

	"eight/internal/models"
)

type HandlerBooks struct {
	store store
}

func NewService(pqdriver *pgxpool.Pool, db *sql.DB) (*HandlerBooks, error) {
	bookStore, err := newStore(pqdriver, db)
	if err != nil {
		return nil, err
	}

	return &HandlerBooks{
		store: bookStore,
	}, nil
}

func (b *HandlerBooks) AllBooks(ctx context.Context) (models.BookSlice, error) {
	books, err := b.store.All(ctx)
	if err != nil {
		return nil, err
	}

	return books, nil
}

func (b *HandlerBooks) CreateBook(ctx context.Context, book *models.Book) (*models.Book, error) {
	return b.store.CreateBook(ctx, book)
}

func (b *HandlerBooks) GetBook(ctx context.Context, bookID int64) (*models.Book, error) {
	return b.store.GetBook(ctx, bookID)
}

func (b *HandlerBooks) Delete(ctx context.Context, bookID int64) error {
	return b.store.Delete(ctx, bookID)
}

func (b *HandlerBooks) Ping() error {
	return b.store.Ping()
}
