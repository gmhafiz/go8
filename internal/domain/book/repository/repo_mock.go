// Code generated by mirip; DO NOT EDIT.
// github.com/gmhafiz/mirip

package repository

import (
	"context"
	"github.com/gmhafiz/go8/internal/domain/book"
)

// BookMock is a mock implementation of Book.
type BookMock struct {
	CreateFunc func(ctx context.Context, bookMiripParam *book.CreateRequest) (uint64, error)
	DeleteFunc func(ctx context.Context, bookID uint64) error
	ListFunc   func(ctx context.Context, f *book.Filter) ([]*book.Schema, error)
	ReadFunc   func(ctx context.Context, bookID uint64) (*book.Schema, error)
	SearchFunc func(ctx context.Context, req *book.Filter) ([]*book.Schema, error)
	UpdateFunc func(ctx context.Context, bookMiripParam *book.UpdateRequest) error
}

func (m *BookMock) Create(ctx context.Context, bookMiripParam *book.CreateRequest) (uint64, error) {
	return m.CreateFunc(ctx, bookMiripParam)
}

func (m *BookMock) Delete(ctx context.Context, bookID uint64) error {
	return m.DeleteFunc(ctx, bookID)
}

func (m *BookMock) List(ctx context.Context, f *book.Filter) ([]*book.Schema, error) {
	return m.ListFunc(ctx, f)
}

func (m *BookMock) Read(ctx context.Context, bookID uint64) (*book.Schema, error) {
	return m.ReadFunc(ctx, bookID)
}

func (m *BookMock) Search(ctx context.Context, req *book.Filter) ([]*book.Schema, error) {
	return m.SearchFunc(ctx, req)
}

func (m *BookMock) Update(ctx context.Context, bookMiripParam *book.UpdateRequest) error {
	return m.UpdateFunc(ctx, bookMiripParam)
}
