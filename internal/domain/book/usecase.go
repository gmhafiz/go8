package book

import (
	"context"

	"github.com/gmhafiz/go8/internal/model"
)

type UseCase interface {
	Create(ctx context.Context, title, description, imageURL, publishedDate string) (*model.Book, error)
	All(ctx context.Context) ([]*model.Book, error)
	Find(ctx context.Context, bookID int64) (*model.Book, error)
	Update(ctx context.Context, book *model.Book) (*model.Book, error)
	Delete(ctx context.Context, bookID int64) error
}
