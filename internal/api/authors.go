package api

import (
	"context"
	"eight/internal/models"
)

func (a API) GetAllAuthors(ctx context.Context) (models.AuthorSlice, error) {
	return a.authors.AllAuthors(ctx)
}

func (a API) CreateAuthor(ctx context.Context, author *models.Author) (*models.Author, error) {
	return a.authors.CreateAuthor(ctx, author)
}

func (a API) GetAuthor(ctx context.Context, authorID int64) (*models.Author, error) {
	return a.authors.GetAuthor(ctx, authorID)
}
