package books

import (
	"context"
	"database/sql"

	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"

	"go8ddd/internal/middleware"
	"go8ddd/internal/model"
)

type BookRepository interface {
	All(ctx context.Context) (model.BookSlice, error)
	Create(ctx context.Context, book *model.Book) (*model.Book, error)
	Get(ctx context.Context, bookID int64) (*model.Book, error)
	Update(ctx context.Context, book *model.Book) (*model.Book, error)
	Delete(ctx context.Context, bookID int64) error
	HardDelete(ctx context.Context, bookID int64) error
}

type repo struct {
	log   zerolog.Logger
	db    *sql.DB
	cache Store
}

func NewRepository(log zerolog.Logger, db *sql.DB, cache *redis.Client) *repo {
	cacheStore, err := newCacheStore(cache, log)
	if err != nil {
		log.Fatal()
	}

	return &repo{
		log:   log,
		db:    db,
		cache: cacheStore,
	}
}

func (r repo) All(ctx context.Context) (model.BookSlice, error) {
	page := ctx.Value("pagination").(middleware.Pagination).Page
	size := ctx.Value("pagination").(middleware.Pagination).Size

	var err error
	var books model.BookSlice

	books = r.cache.All(ctx, page, size)

	if len(books) > 0 {
		return books, nil
	}

	if page == 0 && size == 0 {
		books, err = model.Books().All(ctx, r.db)
	} else {
		books, err = model.Books(
			qm.OrderBy(model.BookColumns.CreatedAt+" DESC"),
			qm.Limit(size),
			qm.Offset(size*(page-1))).
			All(ctx, r.db)
	}

	if err != nil {
		r.log.Error().Msg(err.Error())
		return nil, err
	}

	_ = r.cache.Set(ctx, page, size, &books)

	return books, nil
}

func (r repo) Create(ctx context.Context, book *model.Book) (*model.Book, error) {
	err := book.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		r.log.Error().Msg(err.Error())
		return nil, errors.Wrap(err, err.Error())
	}
	return book, nil
}

func (r repo) Get(ctx context.Context, bookID int64) (*model.Book, error) {
	book, err := model.FindBook(ctx, r.db, bookID)
	if err != nil {
		r.log.Error().Msg(err.Error())
		return nil, errors.Wrap(err, err.Error())
	}
	return book, nil
}

func (r repo) Update(ctx context.Context, book *model.Book) (*model.Book, error) {
	id := ctx.Value("id").(int64)

	bookDB, err := model.FindBook(ctx, r.db, id)
	if err != nil {
		r.log.Error().Msg(err.Error())
		return nil, errors.Wrap(err, err.Error())
	}
	bookDB = book
	_, err = bookDB.Update(ctx, r.db, boil.Infer())
	if err != nil {
		r.log.Error().Msg(err.Error())
		return nil, errors.Wrap(err, err.Error())
	}

	return bookDB, nil
}

func (r repo) Delete(ctx context.Context, bookID int64) error {
	book, err := model.FindBook(ctx, r.db, bookID)
	if err != nil {
		r.log.Error().Msg(err.Error())
		return errors.Wrap(err, err.Error())
	}
	_, err = book.Delete(ctx, r.db, false)
	if err != nil {
		r.log.Error().Msg(err.Error())
		return errors.Wrap(err, err.Error())
	}
	return nil
}

func (r repo) HardDelete(ctx context.Context, bookID int64) error {
	book, err := model.FindBook(ctx, r.db, bookID)
	if err != nil {
		r.log.Error().Msg(err.Error())
		return errors.Wrap(err, err.Error())
	}
	_, err = book.Delete(ctx, r.db, true)
	if err != nil {
		r.log.Error().Msg(err.Error())
		return errors.Wrap(err, err.Error())
	}
	return nil
}
