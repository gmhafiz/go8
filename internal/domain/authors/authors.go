package authors

import (
	"context"
	"database/sql"
	"eight/pkg/elasticsearch"

	"github.com/go-redis/redis/v8"
	"github.com/rs/zerolog"

	"eight/internal/models"
)

type HandlerAuthors struct {
	store store
	cache authorCacheStore
	es    *elasticsearch.Es
}

func NewService(db *sql.DB, logger zerolog.Logger, rdb *redis.Client, es *elasticsearch.Es) (*HandlerAuthors, error) {
	authorStore, err := newStore(db, logger)
	if err != nil {
		return nil, err
	}

	cacheStore, err := newCacheStore(rdb, logger)

	return &HandlerAuthors{
		store: authorStore,
		cache: cacheStore,
		es: es,
	}, nil
}

func (a *HandlerAuthors) AllAuthors(ctx context.Context) (models.AuthorSlice, error) {
	authorRedis, err := a.cache.GetAuthors(ctx)
	if err == nil {
		return authorRedis, nil
	}
	authors, err := a.store.All(ctx)
	if err != nil {
		return nil, err
	}

	_ = a.cache.SetAuthors(ctx, &authors)

	return authors, nil
}

func (a *HandlerAuthors) CreateAuthor(ctx context.Context, author *models.Author) (*models.Author, error) {
	return a.store.CreateAuthor(ctx, author)
}

func (a *HandlerAuthors) GetAuthor(ctx context.Context, authorID int64) (*models.Author, error) {
	return a.store.GetAuthor(ctx, authorID)
}
