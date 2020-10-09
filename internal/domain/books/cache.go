package books

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"
	"github.com/vmihailenco/msgpack/v4"

	"go8ddd/internal/model"
)

type Store interface {
	All(context.Context, int, int) (books model.BookSlice, err error)
	Set(context.Context, int, int, *model.BookSlice) error
}

type cache struct {
	cache  *redis.Client
	logger zerolog.Logger
}

func newCacheStore(redis *redis.Client, logger zerolog.Logger) (*cache, error) {
	return &cache{
		cache:  redis,
		logger: logger,
	}, nil
}

func (c *cache) All(ctx context.Context, page int, size int) (books model.BookSlice,
	err error) {
	var key string
	if page != 0 && size != 0 {
		key = fmt.Sprintf("booksAll-%s-%s", strconv.Itoa(page), strconv.Itoa(size))
	} else {
		key = "booksAll"
	}

	b, err := c.cache.Get(ctx, key).Bytes()
	if err != nil {
		c.logger.Error().Msg(err.Error())
		return nil, err
	}

	err = msgpack.Unmarshal(b, &books)
	if err != nil {
		return nil, err
	}

	return books, nil
}

func (c *cache) Set(ctx context.Context, page, size int, books *model.BookSlice) error {
	var key string
	if page != 0 && size != 0 {
		key = fmt.Sprintf("booksAll-%s-%s", strconv.Itoa(page), strconv.Itoa(size))
	} else {
		key = "booksAll"
	}

	b, err := msgpack.Marshal(books)
	if err != nil {
		c.logger.Error().Msg(err.Error())
		return err
	}

	return c.cache.Set(ctx, key, b, time.Minute*1).Err()
}
