package authors

import (
	"context"

	"go8ddd/internal/model"
)

type AuthorUseCase interface {
	All(ctx context.Context) (model.AuthorSlice, error)
}

type useCase struct {
	authorRepo AuthorRepository
}

func (u *useCase) All(ctx context.Context) (model.AuthorSlice, error) {
	list, err := u.authorRepo.All(ctx)
	if err != nil {
		return nil, err
	}
	return list, nil
}

func NewUseCase(repo AuthorRepository) AuthorUseCase {
	return &useCase{
		authorRepo: repo,
	}
}
