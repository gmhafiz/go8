package authors

import (
	"context"
	"database/sql"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"

	"eight/internal/models"
	"eight/pkg/jobs"
)

type HandlerAuthors struct {
	store store
	cache authorCacheStore
	pool  *jobs.Jobs
}

func NewService(db *sql.DB, logger zerolog.Logger, rdb *redis.Client, pool *jobs.Jobs) (*HandlerAuthors, error) {
	authorStore, err := newStore(db, logger)
	if err != nil {
		return nil, err
	}

	cacheStore, err := newCacheStore(rdb, logger)

	return &HandlerAuthors{
		store: authorStore,
		cache: cacheStore,
		pool: pool,
	}, nil
}

func (h *HandlerAuthors) AllAuthors(ctx context.Context) (models.AuthorSlice, error) {
	authorRedis, err := h.cache.GetAuthors(ctx)
	if err == nil {
		return authorRedis, nil
	}
	authors, err := h.store.All(ctx)
	if err != nil {
		return nil, err
	}

	_ = h.cache.SetAuthors(ctx, &authors)

	return authors, nil
}

func (h *HandlerAuthors) CreateAuthor(ctx context.Context, author *models.Author) (*models.Author, error) {
	return h.store.CreateAuthor(ctx, author)
}

func (h *HandlerAuthors) GetAuthor(ctx context.Context, authorID int64) (*models.Author, error) {
	return h.store.GetAuthor(ctx, authorID)
}
