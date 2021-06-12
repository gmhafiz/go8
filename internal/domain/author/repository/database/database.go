package database

import (
	"context"
	"github.com/gmhafiz/go8/internal/domain/book"

	"github.com/jmoiron/sqlx"

	"github.com/gmhafiz/go8/internal/domain/author"
	"github.com/gmhafiz/go8/internal/models"
)

type repository struct {
	db *sqlx.DB
}

func (r *repository) ReadWithBooks(ctx context.Context, id uint64) (*models.Author, error) {
	a, err := r.Read(ctx, id)
	if err != nil {
		return nil, err
	}

	var books models.BookSlice
	rows, err := r.db.QueryContext(ctx, `SELECT b.*
FROM books b
         INNER JOIN book_authors ba on ba.books_id = b.book_id
         INNER JOIN authors a on a.author_id = ba.author_id
where a.author_id = $1`, a.AuthorID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var b models.Book
		err = rows.Scan(&b.BookID, &b.Title, &b.PublishedDate, &b.ImageURL, &b.Description, &b.CreatedAt, &b.UpdatedAt, &b.DeletedAt)
		books = append(books, &b)
	}

	a.Books = books
	return a, nil
}

func New(db *sqlx.DB) *repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, author *models.Author) (uint64, error) {
	panic("implement me")
}

func (r *repository) List(ctx context.Context, f *author.Filter) ([]*models.Author, error) {
	var authors []*models.Author
	err := r.db.SelectContext(ctx, &authors, "SELECT * from authors")
	if err != nil {
		return nil, err
	}
	return authors, nil
}

func (r *repository) Read(ctx context.Context, authorID uint64) (*models.Author, error) {
	var a models.Author
	err := r.db.GetContext(ctx, &a, "SELECT * from authors where author_id = $1", authorID)
	if err != nil {
		return nil, err
	}

	return &a, nil
}

func (r *repository) Update(ctx context.Context, author *models.Author) error {
	panic("implement me")
}

func (r *repository) Delete(ctx context.Context, authorID uint64) error {
	panic("implement me")
}

type authorbook struct{}

func (r authorbook) Create(ctx context.Context, book *models.Book) (int64, error) {
	panic("implement me")
}

func (r authorbook) List(ctx context.Context, f *book.Filter) ([]*models.Book, error) {
	panic("implement me")
}

func (r authorbook) Read(ctx context.Context, bookID int64) (*models.Book, error) {
	panic("implement me")
}

func (r authorbook) Update(ctx context.Context, book *models.Book) error {
	panic("implement me")
}

func (r authorbook) Delete(ctx context.Context, bookID int64) error {
	panic("implement me")
}

func (r authorbook) Search(ctx context.Context, req *book.Filter) ([]*models.Book, error) {
	panic("implement me")
}
