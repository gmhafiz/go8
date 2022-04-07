package usecase

import (
	"context"
	"errors"

	"github.com/gmhafiz/go8/ent/gen"
	"github.com/gmhafiz/go8/internal/domain/author"
	"github.com/gmhafiz/go8/internal/domain/author/repository"
)

type AuthorUseCase struct {
	repo repository.Author

	searchRepo repository.Searcher

	cacheLRU   repository.AuthorLRUService
	cacheRedis repository.AuthorRedisService
}

//go:generate mirip -rm -out usecase_mock.go . Author
type Author interface {
	Create(ctx context.Context, a *author.CreateRequest) (*gen.Author, error)
	List(ctx context.Context, f *author.Filter) ([]*gen.Author, int, error)
	Read(ctx context.Context, authorID uint) (*gen.Author, error)
	Update(ctx context.Context, author *author.Update) (*gen.Author, error)
	Delete(ctx context.Context, authorID uint) error
}

func New(repo repository.Author, searcher repository.Searcher, cache repository.AuthorLRUService, redisCache repository.AuthorRedisService) *AuthorUseCase {
	return &AuthorUseCase{
		repo:       repo,
		searchRepo: searcher,
		cacheLRU:   cache,
		cacheRedis: redisCache,
	}
}

func (u *AuthorUseCase) Create(ctx context.Context, r *author.CreateRequest) (*gen.Author, error) {
	return u.repo.Create(ctx, r)
}

func (u *AuthorUseCase) List(ctx context.Context, f *author.Filter) ([]*gen.Author, int, error) {
	if f.Base.Search {
		return u.searchRepo.Search(ctx, f)
	}
	// Use cacheRedis layer which is faster. But depends on TTL for cache
	// invalidation. It will call repository layer if cache key is not found.
	return u.cacheRedis.List(ctx, f)

	// Call the database layer.
	//return u.repo.List(ctx, f)
}

func (u *AuthorUseCase) Read(ctx context.Context, authorID uint) (*gen.Author, error) {
	if authorID == 0 {
		return nil, errors.New("ID cannot be 0")
	}
	return u.repo.Read(ctx, authorID)
}

func (u *AuthorUseCase) Update(ctx context.Context, author *author.Update) (*gen.Author, error) {
	// Call cache layer instead to invalidate cache
	return u.cacheRedis.Update(ctx, author)
}

func (u *AuthorUseCase) Delete(ctx context.Context, authorID uint) error {
	if authorID <= 0 {
		return errors.New("ID cannot be 0 or less")
	}
	// As above
	return u.cacheRedis.Delete(ctx, authorID)
}
