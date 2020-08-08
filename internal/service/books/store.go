package books

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"github.com/volatiletech/sqlboiler/v4/boil"

	"eight/internal/models"
)

type store interface {
	All(context.Context) (models.BookSlice, error)
	CreateBook(ctx context.Context, bookID *models.Book) (*models.Book, error)
	GetBook(context.Context, int64) (*models.Book, error)
	Delete(ctx context.Context, bookID int64) error
	Ping() error
}

type bookStore struct {
	db    *sql.DB
}

func (bs *bookStore) All(ctx context.Context) (models.BookSlice, error) {
	bookSlice, err := models.Books().All(ctx, bs.db)
	if err != nil {
		return nil, err
	}

	return bookSlice, nil
}

func (bs *bookStore) CreateBook(ctx context.Context, book *models.Book) (*models.Book, error) {
	//boil.DebugMode = true
	err := book.Insert(ctx, bs.db, boil.Infer())
	if err != nil {
		return book, err
	}
	return book, nil
}

func (bs *bookStore) GetBook(ctx context.Context, bookID int64) (*models.Book, error) {
	var b *models.Book

	book, err := models.Books(models.BookWhere.BookID.EQ(bookID)).One(ctx, bs.db)
	if err != nil {
		return b, errors.Wrap(err, "book not found")
	}

	return book, nil
}

func (bs *bookStore) Delete(ctx context.Context, bookID int64) error {
	book, err := models.FindBook(ctx, bs.db, bookID)
	if err != nil {
		return err
	}
	_, err = book.Delete(ctx, bs.db)
	if err != nil {
		return err
	}

	return nil
}

func (bs *bookStore) Ping() error {
	return bs.db.Ping()
}

func newStore(db *sql.DB) (*bookStore, error) {
	return &bookStore{
		db:    db,
	}, nil
}
