package authors

import (
	"context"
	"eight/internal/middleware"
	"fmt"
	"github.com/vmihailenco/msgpack/v4"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"

	"eight/internal/models"
)

type authorCacheStore interface {
	GetAuthors(ctx context.Context) (authors models.AuthorSlice, err error)
	SetAuthors(ctx context.Context, authors *models.AuthorSlice) error
}

type authorCache struct {
	cache  *redis.Client
	logger zerolog.Logger
}

func newCacheStore(cache *redis.Client, logger zerolog.Logger) (*authorCache, error) {
	return &authorCache{
		cache:  cache,
		logger: logger,
	}, nil
}

func (cache *authorCache) GetAuthors(ctx context.Context) (authors models.AuthorSlice, err error) {
	from := ctx.Value("pagination").(middleware.Pagination).Page
	size := ctx.Value("pagination").(middleware.Pagination).Size

	var key string
	if from != 0 && size != 0 {
		key = fmt.Sprintf("authorsAll-%s-%s", strconv.Itoa(from), strconv.Itoa(size))
	} else {
		key = "authorsAll"
	}

	a, err := cache.cache.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	err = msgpack.Unmarshal(a, &authors)
	if err != nil {
		cache.logger.Error().Msg(err.Error())
		return nil, err
	}

	return authors, nil
}

func (cache *authorCache) SetAuthors(ctx context.Context, authors *models.AuthorSlice) error {
	from := ctx.Value("pagination").(middleware.Pagination).Page
	size := ctx.Value("pagination").(middleware.Pagination).Size

	var key string
	if from != 0 && size != 0 {
		key = fmt.Sprintf("authorsAll-%s-%s", strconv.Itoa(from), strconv.Itoa(size))
	} else {
		key = "authorsAll"
	}

	a, err := msgpack.Marshal(authors)
	if err != nil {
		cache.logger.Error().Msg(err.Error())
		return err
	}

	return cache.cache.Set(ctx, key, a, time.Minute*1).Err()
}
