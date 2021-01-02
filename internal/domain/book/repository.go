package book

import (
	"context"

	"github.com/gmhafiz/go8/internal/model"
	"github.com/gmhafiz/go8/internal/resource"
)

type Repository interface {
	Create(ctx context.Context, book *model.Book) (int64, error)
	All(ctx context.Context) ([]resource.BookDB, error)
	Find(ctx context.Context, bookID int64) (*model.Book, error)
	Close()
	Drop() error
	Up() error
}
