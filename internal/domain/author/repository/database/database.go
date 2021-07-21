package database

import (
	"context"
	"fmt"
	"github.com/jmoiron/sqlx"

	"github.com/gmhafiz/go8/internal/domain/author"
	"github.com/gmhafiz/go8/internal/models"
)

type repository struct {
	db *sqlx.DB
}

func (r *repository) ReadWithBooks(ctx context.Context, id uint64) (*author.AuthorB, error) {
	a, err := r.Read(ctx, id)
	if err != nil {
		return nil, err
	}

	var books models.BookSlice
	rows, err := r.db.QueryContext(ctx, `SELECT b.book_id
FROM books b
         INNER JOIN book_authors ba on ba.books_id = b.book_id
         INNER JOIN authors a on a.author_id = ba.author_id
where a.author_id = $1`, a.AuthorID)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var b models.Book
		err = rows.Scan(&b.BookID)
		if err != nil {
			return nil, fmt.Errorf("error scanning book")
		}
		books = append(books, &b)
	}

	author := &author.AuthorB{
		Author: a,
		Books:  books,
	}

	return author, nil
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
