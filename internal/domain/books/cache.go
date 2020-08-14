package books

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"github.com/vmihailenco/msgpack/v4"

	"eight/internal/middleware"
	"eight/internal/models"
)

type bookCacheStore interface {
	GetBooks(ctx context.Context) (books models.BookSlice, err error)
	SetBooks(ctx context.Context, books *models.BookSlice) error
}

type bookCache struct {
	cache *redis.Client
	logger zerolog.Logger
}

func newCacheStore(cache *redis.Client, logger zerolog.Logger) (*bookCache, error) {
	return &bookCache{
		cache: cache,
		logger: logger,
	}, nil
}


func (cache *bookCache) GetBooks(ctx context.Context) (books models.BookSlice, err error) {
	from := ctx.Value("pagination").(middleware.Pagination).Page
	size := ctx.Value("pagination").(middleware.Pagination).Size

	var key string
	if from != 0 && size != 0 {
		key = fmt.Sprintf("booksAll-%s-%s", strconv.Itoa(from), strconv.Itoa(size))
	} else {
		key = "booksAll"
	}

	b, err := cache.cache.Get(ctx, key).Bytes()
	if err != nil {
		cache.logger.Error().Msg(err.Error())
		return nil, err
	}

	err = msgpack.Unmarshal(b, &books)
	if err != nil {
		cache.logger.Error().Msg(err.Error())
		return nil, err
	}

	return books, nil
}

func (cache *bookCache) SetBooks(ctx context.Context, books *models.BookSlice) error {
	from := ctx.Value("pagination").(middleware.Pagination).Page
	size := ctx.Value("pagination").(middleware.Pagination).Size

	var key string
	if from != 0 && size != 0 {
		key = fmt.Sprintf("booksAll-%s-%s", strconv.Itoa(from), strconv.Itoa(size))
	} else {
		key = "booksAll"
	}

	b, err := msgpack.Marshal(books)
	if err != nil {
		cache.logger.Error().Msg(err.Error())
		return err
	}

	return cache.cache.Set(ctx, key, b, time.Minute * 1).Err()
}
