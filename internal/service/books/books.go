package books

import (
	"context"
	"database/sql"
	"github.com/go-redis/redis/v8"

	"eight/internal/models"
)

type HandlerBooks struct {
	store store
	cache bookCacheStore
}

func NewService(db *sql.DB, rdb *redis.Client) (*HandlerBooks, error) {
	bookStore, err := newStore(db)
	if err != nil {
		return nil, err
	}

	cacheStore, err := newCacheStore(rdb)

	return &HandlerBooks{
		store: bookStore,
		cache: cacheStore,
	}, nil
}

func (b *HandlerBooks) AllBooks(ctx context.Context) (models.BookSlice, error) {
	booksRedis, err := b.cache.GetBooks(ctx)
	if err == nil {
		return booksRedis, nil
	}
	books, err := b.store.All(ctx)
	if err != nil {
		return nil, err
	}

	_ = b.cache.SetBooks(ctx, &books)

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
