package book

import (
	"context"

	"github.com/gmhafiz/go8/internal/models"
)

type Repository interface {
	Create(ctx context.Context, book *models.Book) (int64, error)
	List(ctx context.Context, f *Filter) ([]*models.Book, error)
	Read(ctx context.Context, bookID int64) (*models.Book, error)
	Update(ctx context.Context, book *models.Book) error
	Delete(ctx context.Context, bookID int64) error
	Search(ctx context.Context, req *Filter) ([]*models.Book, error)
}
