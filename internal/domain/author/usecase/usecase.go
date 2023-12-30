package usecase

import (
	"context"
	"errors"

	"go.opentelemetry.io/otel"

	"github.com/gmhafiz/go8/config"
	"github.com/gmhafiz/go8/internal/domain/author"
	"github.com/gmhafiz/go8/internal/domain/author/repository"
)

type AuthorUseCase struct {
	repo repository.Author

	searchRepo repository.Searcher

	cacheLRU   repository.AuthorLRUService
	cacheRedis repository.AuthorRedisService
	cfg        config.Cache
}

//go:generate mirip -rm -out usecase_mock.go . Author
type Author interface {
	Create(ctx context.Context, a *author.CreateRequest) (*author.Schema, error)
	List(ctx context.Context, f *author.Filter) ([]*author.Schema, int, error)
	Read(ctx context.Context, authorID uint64) (*author.Schema, error)
	Update(ctx context.Context, author *author.UpdateRequest) (*author.Schema, error)
	Delete(ctx context.Context, authorID uint64) error
}

func New(c config.Cache, repo repository.Author, searcher repository.Searcher, cache repository.AuthorLRUService, redisCache repository.AuthorRedisService) *AuthorUseCase {
	return &AuthorUseCase{
		cfg:        c,
		repo:       repo,
		searchRepo: searcher,
		cacheLRU:   cache,
		cacheRedis: redisCache,
	}
}

func (u *AuthorUseCase) Create(ctx context.Context, r *author.CreateRequest) (*author.Schema, error) {
	return u.repo.Create(ctx, r)
}

func (u *AuthorUseCase) List(ctx context.Context, f *author.Filter) ([]*author.Schema, int, error) {
	tracer := otel.Tracer("")
	ctx, span := tracer.Start(ctx, "AuthorUseCase")
	defer span.End()

	if f.Base.Search {
		return u.searchRepo.Search(ctx, f)
	}

	if u.cfg.Enable {
		// Use cacheRedis layer which is faster. But depends on TTL for cache
		// invalidation. It will call repository layer if cache key is not found.
		return u.cacheRedis.List(ctx, f)
	}

	// Call the database layer.
	return u.repo.List(ctx, f)
}

func (u *AuthorUseCase) Read(ctx context.Context, authorID uint64) (*author.Schema, error) {
	if authorID == 0 {
		return nil, errors.New("ID cannot be 0")
	}
	return u.repo.Read(ctx, authorID)
}

func (u *AuthorUseCase) Update(ctx context.Context, author *author.UpdateRequest) (*author.Schema, error) {
	if u.cfg.Enable {
		// Call cache layer instead to invalidate cache
		return u.cacheRedis.Update(ctx, author)
	}

	return u.repo.Update(ctx, author)
}

func (u *AuthorUseCase) Delete(ctx context.Context, authorID uint64) error {
	if authorID <= 0 {
		return errors.New("ID cannot be 0 or less")
	}

	if u.cfg.Enable {
		// As above
		return u.cacheRedis.Delete(ctx, authorID)
	}

	return u.repo.Delete(ctx, authorID)
}
