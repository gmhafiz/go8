package authors

import (
	"context"

	"go8ddd/internal/model"
)

type AuthorUseCase interface {
	All(ctx context.Context) (model.AuthorSlice, error)
}

type authorUseCase struct {
	authorRepo AuthorRepository
}

func (u *authorUseCase) All(ctx context.Context) (model.AuthorSlice, error) {
	list, err := u.authorRepo.All(ctx)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func NewUseCase(repo AuthorRepository) AuthorUseCase {
	return &authorUseCase{
		authorRepo: repo,
	}
}
