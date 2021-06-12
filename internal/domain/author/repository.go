package author

import (
	"context"

	"github.com/gmhafiz/go8/internal/models"
)

type Repository interface {
	Create(ctx context.Context, Author *models.Author) (uint64, error)
	List(ctx context.Context, f *Filter) ([]*models.Author, error)
	Read(ctx context.Context, authorID uint64) (*models.Author, error)
	Update(ctx context.Context, author *models.Author) error
	Delete(ctx context.Context, authorID uint64) error
	ReadWithBooks(ctx context.Context, id uint64) (*models.Author, error)
}









