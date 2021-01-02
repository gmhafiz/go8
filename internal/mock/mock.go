package mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/gmhafiz/go8/internal/model"
)

type BookUseCaseMock struct {
	mock.Mock
}

func (m *BookUseCaseMock) All(context.Context) ([]*model.Book, error) {
	args := m.Called()
	return args.Get(0).([]*model.Book), args.Error(1)
}

func (m *BookUseCaseMock) Create(ctx context.Context, title, description, imageURL, publishedDate string) (*model.Book, error) {
	args := m.Called(title, description, imageURL, publishedDate)
	return args.Get(0).(*model.Book), args.Error(1)
}

func (m *BookUseCaseMock) Find(ctx context.Context, bookID int64) (*model.Book, error) {
	panic("implement me")
}
