package mocks

import (
	"context"

	"github.com/stretchr/testify/mock"

	"go8ddd/internal/model"
)

type BookUseCase struct {
	mock.Mock
}

func (_m *BookUseCase) All(ctx context.Context) ([]model.Book, error) {

	ret := _m.Called(ctx)

	var r0 []model.Book
	if rf, ok := ret.Get(0).(func(context.Context) []model.Book); ok {
		r0 = rf(ctx)
	}

	return r0, nil
}
