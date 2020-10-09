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

type BookRepo struct {
	Log   zerolog.Logger
	DB    *sql.DB
	Cache Store
}

func NewRepository(log zerolog.Logger, db *sql.DB, cache *redis.Client) *BookRepo {
	cacheStore, err := newCacheStore(cache, log)
	if err != nil {
		log.Fatal()
	}

	return &BookRepo{
		Log:   log,
		DB:    db,
		Cache: cacheStore,
	}
}

func (r BookRepo) All(ctx context.Context) (model.BookSlice, error) {
	page := ctx.Value("pagination").(middleware.Pagination).Page
	size := ctx.Value("pagination").(middleware.Pagination).Size

	var err error
	var books model.BookSlice

	books, err = r.Cache.All(ctx, page, size)
	if err != nil {
		r.Log.Warn().Msg(err.Error())
	}

	if len(books) > 0 {
		return books, nil
	}

	if page == 0 && size == 0 {
		books, err = model.Books().All(ctx, r.DB)
	} else {
		books, err = model.Books(
			qm.OrderBy(model.BookColumns.CreatedAt+" DESC"),
			qm.Limit(size),
			qm.Offset(size*(page-1))).
			All(ctx, r.DB)
	}

	if err != nil {
		r.Log.Error().Msg(err.Error())
		return nil, err
	}

	err = r.Cache.Set(ctx, page, size, &books)
	if err != nil {
		return nil, err
	}

	return books, nil
}

func (r BookRepo) Create(ctx context.Context, book *model.Book) (*model.Book, error) {
	err := book.Insert(ctx, r.DB, boil.Infer())
	if err != nil {
		r.Log.Error().Msg(err.Error())
		return nil, errors.Wrap(err, err.Error())
	}
	return book, nil
}

func (r BookRepo) Get(ctx context.Context, bookID int64) (*model.Book, error) {
	book, err := model.FindBook(ctx, r.DB, bookID)
	if err != nil {
		r.Log.Error().Msg(err.Error())
		return nil, errors.Wrap(err, err.Error())
	}
	return book, nil
}

func (r BookRepo) Update(ctx context.Context, book *model.Book) (*model.Book, error) {
	id := ctx.Value("id").(int64)

	bookDB, err := model.FindBook(ctx, r.DB, id)
	if err != nil {
		r.Log.Error().Msg(err.Error())
		return nil, errors.Wrap(err, err.Error())
	}
	bookDB = book
	_, err = bookDB.Update(ctx, r.DB, boil.Infer())
	if err != nil {
		r.Log.Error().Msg(err.Error())
		return nil, errors.Wrap(err, err.Error())
	}

	return bookDB, nil
}

func (r BookRepo) Delete(ctx context.Context, bookID int64) error {
	book, err := model.FindBook(ctx, r.DB, bookID)
	if err != nil {
		r.Log.Error().Msg(err.Error())
		return errors.Wrap(err, err.Error())
	}
	_, err = book.Delete(ctx, r.DB, false)
	if err != nil {
		r.Log.Error().Msg(err.Error())
		return errors.Wrap(err, err.Error())
	}
	return nil
}

func (r BookRepo) HardDelete(ctx context.Context, bookID int64) error {
	book, err := model.FindBook(ctx, r.DB, bookID)
	if err != nil {
		r.Log.Error().Msg(err.Error())
		return errors.Wrap(err, err.Error())
	}
	_, err = book.Delete(ctx, r.DB, true)
	if err != nil {
		r.Log.Error().Msg(err.Error())
		return errors.Wrap(err, err.Error())
	}
	return nil
}
