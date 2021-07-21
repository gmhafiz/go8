package usecase

import (
	"context"
	"github.com/gmhafiz/go8/internal/domain/book"

	"github.com/gmhafiz/go8/internal/domain/author"
	"github.com/gmhafiz/go8/internal/models"
)

type useCase struct {
	repo     author.Repository
	bookRepo book.Repository
}

func (u *useCase) ReadWithBooks(ctx context.Context, authorID uint64) (*author.AuthorB, error) {
	books, err := u.repo.ReadWithBooks(ctx, authorID)
	if err != nil {
		return nil, err
	}

	return books, nil
}

func New(repo author.Repository, bookRepo book.Repository) *useCase {
	return &useCase{
		repo:     repo,
		bookRepo: bookRepo,
	}
}

func (u *useCase) Create(ctx context.Context, r author.Request) (*models.Author, error) {
	bk := author.ToAuthor(&r)
	authorID, err := u.repo.Create(ctx, bk)
	if err != nil {
		return nil, err
	}
	authorFound, err := u.repo.Read(ctx, authorID)
	if err != nil {
		return nil, err
	}
	return authorFound, err
}

func (u *useCase) List(ctx context.Context, f *author.Filter) ([]*models.Author, error) {
	return u.repo.List(ctx, f)
}

func (u *useCase) Read(ctx context.Context, authorID uint64) (*models.Author, error) {
	return u.repo.Read(ctx, authorID)
}

func (u *useCase) Update(ctx context.Context, author *models.Author) (*models.Author, error) {
	err := u.repo.Update(ctx, author)
	if err != nil {
		return nil, err
	}
	return u.repo.Read(ctx, uint64(author.AuthorID))
}

func (u *useCase) Delete(ctx context.Context, authorID uint64) error {
	return u.repo.Delete(ctx, authorID)
}
