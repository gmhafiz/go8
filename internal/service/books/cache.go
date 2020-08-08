package books

import (
	"context"

	"github.com/go-redis/redis/v8"
	"github.com/vmihailenco/msgpack/v4"

	"eight/internal/models"
)

type bookCacheStore interface {
	GetBooks(ctx context.Context) (books models.BookSlice, err error)
	SetBooks(ctx context.Context, books *models.BookSlice) error
}

type bookCache struct {
	cache *redis.Client
}

func newCacheStore(cache *redis.Client) (*bookCache, error) {
	return &bookCache{
		cache: cache,
	}, nil
}


func (cache *bookCache) GetBooks(ctx context.Context) (books models.BookSlice, err error) {
	b, err := cache.cache.Get(ctx, "booksAll").Bytes()
	if err != nil {
		return nil, err
	}

	err = msgpack.Unmarshal(b, &books)
	if err != nil {
		return nil, err
	}

	return books, nil
}

func (cache *bookCache) SetBooks(ctx context.Context, books *models.BookSlice) error {
	b, err := msgpack.Marshal(books)
	if err != nil {
		return err
	}

	err = cache.cache.Set(ctx, "booksAll", b, 0).Err()
	return err
}
