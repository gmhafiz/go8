package book

import (
	"context"

	"github.com/gmhafiz/go8/internal/model"
	"github.com/gmhafiz/go8/internal/resource"
)

type UseCase interface {
	Create(ctx context.Context, title, description, imageURL, publishedDate string) (*model.Book, error)
	All(ctx context.Context) ([]resource.BookDB, error)
	Find(ctx context.Context, bookID int64) (*model.Book, error)
}
