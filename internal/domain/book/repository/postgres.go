package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"

	"github.com/gmhafiz/go8/internal/domain/book"
	"github.com/gmhafiz/go8/internal/utility/message"
)

//go:generate mirip -rm -pkg repository -out repo_mock.go . Book
type Book interface {
	Create(ctx context.Context, book *book.CreateRequest) (int, error)
	List(ctx context.Context, f *book.Filter) ([]*book.Schema, error)
	Read(ctx context.Context, bookID int) (*book.Schema, error)
	Update(ctx context.Context, book *book.UpdateRequest) error
	Delete(ctx context.Context, bookID int) error
	Search(ctx context.Context, req *book.Filter) ([]*book.Schema, error)
}

type bookRepository struct {
	db *sqlx.DB
}

const (
	InsertIntoBooks         = "INSERT INTO books (title, published_date, image_url, description) VALUES ($1, $2, $3, $4) RETURNING id"
	SelectFromBooks         = "SELECT * FROM books ORDER BY created_at DESC"
	SelectFromBooksPaginate = "SELECT * FROM books ORDER BY created_at DESC LIMIT $1 OFFSET $2"
	SelectBookByID          = "SELECT * FROM books where id = $1"
	UpdateBook              = "UPDATE books set title = $1, description = $2, published_date = $3, image_url = $4 where id = $5 RETURNING id"
	DeleteByID              = "DELETE FROM books where id = ($1) RETURNING id"
	SearchBooks             = "SELECT * FROM books where title like '%' || $1 || '%' and description like '%'|| $2 || '%' ORDER BY published_date DESC"
	SearchBooksPaginate     = "SELECT * FROM books where title like '%' || '%' || $1 || '%' || '%' and description like '%'|| $2 || '%' ORDER BY published_date DESC LIMIT $3 OFFSET $4"
)

func New(db *sqlx.DB) *bookRepository {
	return &bookRepository{db: db}
}

func (r *bookRepository) Create(ctx context.Context, req *book.CreateRequest) (bookID int, err error) {
	if err = r.db.QueryRowContext(ctx, InsertIntoBooks, req.Title, req.PublishedDate, req.ImageURL, req.Description).Scan(&bookID); err != nil {
		return 0, errors.New("repository.Book.Create")
	}

	return bookID, nil
}

func (r *bookRepository) List(ctx context.Context, f *book.Filter) ([]*book.Schema, error) {
	if f == nil {
		return nil, errors.New("filter cannot be nil")
	}
	if f.Base.DisablePaging {
		var books []*book.Schema
		err := r.db.SelectContext(ctx, &books, SelectFromBooks)
		if err != nil {
			return nil, message.ErrFetchingBook
		}

		return books, nil
	} else {
		var books []*book.Schema
		err := r.db.SelectContext(ctx, &books, SelectFromBooksPaginate, f.Base.Limit, f.Base.Offset)
		if err != nil {
			return nil, message.ErrFetchingBook
		}
		return books, nil
	}
}

func (r *bookRepository) Read(ctx context.Context, bookID int) (*book.Schema, error) {
	var b book.Schema
	err := r.db.GetContext(ctx, &b, SelectBookByID, bookID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, message.ErrBadRequest
		}
		return nil, err
	}

	return &b, err
}

func (r *bookRepository) Update(ctx context.Context, book *book.UpdateRequest) error {
	var returnedID int

	err := r.db.QueryRowContext(ctx, UpdateBook,
		book.Title,
		book.Description,
		book.PublishedDate,
		book.ImageURL,
		book.ID,
	).Scan(&returnedID)
	if err != nil {
		return err
	}

	return nil
}

func (r *bookRepository) Delete(ctx context.Context, bookID int) error {
	var returnedID int
	err := r.db.QueryRowContext(ctx, DeleteByID, bookID).Scan(&returnedID)
	if err != nil {
		return fmt.Errorf("ID not found: %w", err)
	}

	return nil
}

func (r *bookRepository) Search(ctx context.Context, f *book.Filter) ([]*book.Schema, error) {
	if f == nil {
		return nil, errors.New("filter cannot be nil")
	}
	var books []*book.Schema
	err := r.db.SelectContext(ctx, &books, SearchBooksPaginate,
		f.Title,
		f.Description,
		f.Base.Limit,
		f.Base.Offset,
	)
	if err != nil {
		return nil, err
	}

	return books, nil
}
