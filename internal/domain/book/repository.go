package book

import (
	"context"

	"github.com/gmhafiz/go8/internal/model"
)

type Repository interface {
	Create(ctx context.Context, book *model.Book) (int64, error)
	All(ctx context.Context) ([]*model.Book, error)
	Find(ctx context.Context, bookID int64) (*model.Book, error)
	Update(ctx context.Context, book *model.Book) (*model.Book, error)
	Delete(ctx context.Context, bookID int64) error
	Close()
	Drop() error
	Up() error
}
