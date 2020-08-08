package books

import (
	"context"
	"eight/internal/middleware"
	"fmt"

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
	from := ctx.Value("pagination").(middleware.Pagination).Page
	size := ctx.Value("pagination").(middleware.Pagination).Size

	var key string
	if from != 0 && size != 0 {
		key = fmt.Sprintf("booksAll-%s-%s", string(from), string(size))
	} else {
		key = "booksALl"
	}

	b, err := cache.cache.Get(ctx, key).Bytes()
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
	from := ctx.Value("pagination").(middleware.Pagination).Page
	size := ctx.Value("pagination").(middleware.Pagination).Size

	var key string
	if from != 0 && size != 0 {
		key = fmt.Sprintf("booksAll-%s-%s", string(from), string(size))
	} else {
		key = "booksAll"
	}

	b, err := msgpack.Marshal(books)
	if err != nil {
		return err
	}

	return cache.cache.Set(ctx, key, b, 0).Err()
}
