package book

import (
	"context"

	"github.com/gmhafiz/go8/internal/models"
)

type Repository interface {
	Create(ctx context.Context, book *models.Book) (int64, error)
	All(ctx context.Context) ([]*models.Book, error)
	Find(ctx context.Context, bookID int64) (*models.Book, error)
	Update(ctx context.Context, book *models.Book) (*models.Book, error)
	Delete(ctx context.Context, bookID int64) error
	Search(ctx context.Context, req *Request) ([]*models.Book, error)
}

type Test interface {
	Repository
	Close()
	Drop() error
	Up() error
}
