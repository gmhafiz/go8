package mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/gmhafiz/go8/internal/model"
	"github.com/gmhafiz/go8/internal/resource"
)

type BookUseCaseMock struct {
	mock.Mock
}

func (m *BookUseCaseMock) All(context.Context) ([]resource.BookDB, error) {
	args := m.Called()
	return args.Get(0).([]resource.BookDB), args.Error(1)
}

func (m *BookUseCaseMock) Create(ctx context.Context, title, description, imageURL, publishedDate string) (*model.Book, error) {
	args := m.Called(title, description, imageURL, publishedDate)
	return args.Get(0).(*model.Book), args.Error(1)
}

func (m *BookUseCaseMock) Find(ctx context.Context, bookID int64) (*model.Book, error) {
	panic("implement me")
}
