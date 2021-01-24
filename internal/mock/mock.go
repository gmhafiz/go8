package mock

import (
	"context"

	"github.com/stretchr/testify/mock"

	"github.com/gmhafiz/go8/internal/model"
)

type BookUseCaseMock struct {
	mock.Mock
}

func (m *BookUseCaseMock) Update(ctx context.Context, book *model.Book) (*model.Book, error) {
	panic("implement me")
}

func (m *BookUseCaseMock) Delete(ctx context.Context, bookID int64) error {
	panic("implement me")
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
