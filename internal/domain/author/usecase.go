package author

import (
	"context"

	"github.com/gmhafiz/go8/internal/models"
)

type UseCase interface {
	Create(ctx context.Context, author Request) (*models.Author, error)
	List(ctx context.Context, f *Filter) ([]*models.Author, error)
	Read(ctx context.Context, authorID uint64) (*models.Author, error)
	Update(ctx context.Context, author *models.Author) (*models.Author, error)
	Delete(ctx context.Context, authorID uint64) error
	ReadWithBooks(ctx context.Context, u uint64) (*models.Author, error)
}