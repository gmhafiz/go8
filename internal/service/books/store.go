package books

import (
	"context"
	"database/sql"
	"log"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v4/pgxpool"
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
}

type bookStore struct {
	qbuilder squirrel.StatementBuilderType
	pqdriver *pgxpool.Pool
	db       *sql.DB
}

func (bs *bookStore) All(ctx context.Context) (models.BookSlice, error) {
	bookSlice, err := models.Books().All(ctx, bs.db)
	if err != nil {
		log.Println(err)
	}

	return bookSlice, nil
}

func (bs *bookStore) CreateBook(ctx context.Context, book *models.Book) (*models.Book, error) {
	err := book.Insert(ctx, bs.db, boil.Infer())
	if err != nil {
		return book, err
	}
	return book, nil
}

func (bs *bookStore) GetBook(ctx context.Context, bookID int64) (*models.Book, error) {
	boil.DebugMode = true
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

func newStore(pqdriver *pgxpool.Pool, db *sql.DB) (*bookStore, error) {
	return &bookStore{
		qbuilder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
		pqdriver: pqdriver,
		db:       db,
	}, nil
}
