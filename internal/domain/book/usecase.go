package book

import (
	"context"

	"github.com/gmhafiz/go8/internal/model"
	"github.com/gmhafiz/go8/internal/resource"
)

type UseCase interface {
	Create(ctx context.Context, title, description string) (*model.Book, error)
	All(ctx context.Context) ([]resource.BookDB, error)
}
