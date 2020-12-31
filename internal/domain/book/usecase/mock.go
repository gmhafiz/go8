package usecase

import (
	"context"
	"github.com/gmhafiz/go8/internal/model"

	"github.com/stretchr/testify/mock"

	"github.com/gmhafiz/go8/internal/resource"
)

type BookUseCaseMock struct {
	mock.Mock
}

func (m *BookUseCaseMock) All(ctx context.Context) ([]resource.BookDB, error) {
	args := m.Called()
	return args.Get(0).([]resource.BookDB), args.Error(1)
}

func (m *BookUseCaseMock) Create(ctx context.Context, title, description string) (*model.Book, error) {
	args := m.Called(title, description)
	return args.Get(0).(*model.Book), args.Error(1)
}
