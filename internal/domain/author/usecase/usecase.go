package usecase

import (
	"context"

	"github.com/gmhafiz/go8/ent/gen"
	"github.com/gmhafiz/go8/internal/domain/author"
	authorCache "github.com/gmhafiz/go8/internal/domain/author/repository/cache"
	"github.com/gmhafiz/go8/internal/domain/author/repository/database"
)

type AuthorUseCase struct {
	repo database.Repository

	searchRepo database.Searcher

	cacheLRU   authorCache.AuthorLRUService
	cacheRedis authorCache.AuthorRedisService
}

type UseCase interface {
	Create(ctx context.Context, a author.CreateRequest) (*gen.Author, error)
	List(ctx context.Context, f *author.Filter) ([]*gen.Author, int64, error)
	Read(ctx context.Context, authorID uint64) (*gen.Author, error)
	Update(ctx context.Context, author *author.Update) (*gen.Author, error)
	Delete(ctx context.Context, authorID int64) error
}

func New(repo database.Repository, searcher database.Searcher, cache authorCache.AuthorLRUService, redisCache authorCache.AuthorRedisService) *AuthorUseCase {
	return &AuthorUseCase{
		repo:       repo,
		searchRepo: searcher,
		cacheLRU:   cache,
		cacheRedis: redisCache,
	}
}

func (u *AuthorUseCase) Create(ctx context.Context, r author.CreateRequest) (*gen.Author, error) {
	return u.repo.Create(ctx, r)
}

func (u *AuthorUseCase) List(ctx context.Context, f *author.Filter) ([]*gen.Author, int64, error) {
	if f.Base.Search {
		return u.searchRepo.Search(ctx, f)
	}
	// Use cacheRedis layer which is faster. But need to remember to invalidate
	// cache when it is stale.
	return u.cacheRedis.List(ctx, f)

	// Call the database layer.
	//return u.repo.List(ctx, f)
}

func (u *AuthorUseCase) Read(ctx context.Context, authorID uint64) (*gen.Author, error) {
	return u.repo.Read(ctx, authorID)
}

func (u *AuthorUseCase) Update(ctx context.Context, author *author.Update) (*gen.Author, error) {
	// Call cache layer instead to invalidate cache
	return u.cacheRedis.Update(ctx, author)
}

func (u *AuthorUseCase) Delete(ctx context.Context, authorID int64) error {
	// As above
	return u.cacheRedis.Delete(ctx, authorID)
}
