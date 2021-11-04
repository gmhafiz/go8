package usecase

import (
	"context"

	"github.com/gmhafiz/go8/internal/domain/author"
	authorCache "github.com/gmhafiz/go8/internal/domain/author/repository/cache"
	"github.com/gmhafiz/go8/internal/domain/author/repository/database"
	"github.com/gmhafiz/go8/internal/domain/book"
	"github.com/gmhafiz/go8/internal/models"
)

type useCase struct {
	repo       database.Repository
	bookRepo   book.Repository
	searchRepo database.Searcher
	cacheLRU   authorCache.AuthorLRUService
	cacheRedis authorCache.AuthorRedisService
}

type UseCase interface {
	Create(ctx context.Context, a author.CreateRequest) (*author.CreateResponse, error)
	List(ctx context.Context, f *author.Filter) ([]*models.Author, int64, error)
	Read(ctx context.Context, authorID uint64) (*models.Author, error)
	Update(ctx context.Context, a *author.Update) (*models.Author, error)
	Delete(ctx context.Context, authorID int64) error
	ReadWithBooks(ctx context.Context, u uint64) (*author.WithBooks, error)
}

func New(repo database.Repository, searcher database.Searcher, cache authorCache.AuthorLRUService, redisCache authorCache.AuthorRedisService, bookRepo book.Repository) *useCase {
	return &useCase{
		repo:       repo,
		bookRepo:   bookRepo,
		searchRepo: searcher,
		cacheLRU:   cache,
		cacheRedis: redisCache,
	}
}

func (u *useCase) Create(ctx context.Context, r author.CreateRequest) (*author.CreateResponse, error) {
	return u.repo.CreateRead(ctx, r)
}

func (u *useCase) List(ctx context.Context, f *author.Filter) ([]*models.Author, int64, error) {
	if f.Base.Search {
		return u.searchRepo.Search(ctx, f)
	}
	// Use cacheRedis layer which is faster. But need to remember to invalidate
	// cache when it is stale.
	return u.cacheRedis.List(ctx, f)

	// Call the database layer.
	//return u.repo.List(ctx, f)
}

func (u *useCase) Read(ctx context.Context, authorID uint64) (*models.Author, error) {
	return u.repo.Read(ctx, authorID)
}

func (u *useCase) ReadWithBooks(ctx context.Context, authorID uint64) (*author.WithBooks, error) {
	return u.repo.ReadWithBooks(ctx, authorID)
}

func (u *useCase) Update(ctx context.Context, author *author.Update) (*models.Author, error) {
	// Call cache layer instead to invalidate cache
	return u.cacheRedis.Update(ctx, author.UpdateToAuthor(author))
}

func (u *useCase) Delete(ctx context.Context, authorID int64) error {
	// As above
	return u.cacheRedis.Delete(ctx, authorID)
}
