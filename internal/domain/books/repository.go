package books

import (
	"context"
	"database/sql"
	"github.com/go-redis/redis/v8"
	"github.com/olivere/elastic/v7"
	"github.com/pkg/errors"
	"github.com/rs/zerolog"
	"github.com/volatiletech/sqlboiler/v4/boil"
	"github.com/volatiletech/sqlboiler/v4/queries/qm"
	"go8ddd/internal/library/elasticsearch"
	"go8ddd/internal/middleware"
	"go8ddd/internal/model"
	"reflect"
)

type BookRepository interface {
	All(ctx context.Context) (model.BookSlice, error)
	Create(ctx context.Context, book *model.Book) (*model.Book, error)
	Get(ctx context.Context, bookID int64) (*model.Book, error)
	Update(ctx context.Context, book *model.Book) (*model.Book, error)
	Delete(ctx context.Context, bookID int64) error
	HardDelete(ctx context.Context, bookID int64) error
	Search(context.Context, string) ([]model.Book, error)
}

type bookRepo struct {
	log   zerolog.Logger
	db    *sql.DB
	cache Store
	es    *elasticsearch.Es
}

func NewRepository(log zerolog.Logger, db *sql.DB, cache *redis.Client, es *elasticsearch.Es) *bookRepo {
	cacheStore, err := newCacheStore(cache, log)
	if err != nil {
		log.Fatal()
	}

	return &bookRepo{
		log:   log,
		db:    db,
		cache: cacheStore,
		es: es,
	}
}

func (r bookRepo) All(ctx context.Context) (model.BookSlice, error) {
	page := ctx.Value("pagination").(middleware.Pagination).Page
	size := ctx.Value("pagination").(middleware.Pagination).Size

	var err error
	var books model.BookSlice

	books, err = r.cache.All(ctx, page, size)
	if err != nil {
		r.log.Error().Msg(err.Error())
	}

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

func (r bookRepo) Create(ctx context.Context, book *model.Book) (*model.Book, error) {
	err := book.Insert(ctx, r.db, boil.Infer())
	if err != nil {
		r.log.Error().Msg(err.Error())
		return nil, errors.Wrap(err, err.Error())
	}

	_, err = r.es.Client.Index().Index("go8-books").BodyJson(book).Do(context.Background())
	if err != nil {
		return book, err
	}

	return book, nil
}

func (r bookRepo) Get(ctx context.Context, bookID int64) (*model.Book, error) {
	book, err := model.FindBook(ctx, r.db, bookID)
	if err != nil {
		r.log.Error().Msg(err.Error())
		return nil, errors.Wrap(err, err.Error())
	}
	return book, nil
}

func (r bookRepo) Update(ctx context.Context, book *model.Book) (*model.Book, error) {
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

func (r bookRepo) Delete(ctx context.Context, bookID int64) error {
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

func (r bookRepo) HardDelete(ctx context.Context, bookID int64) error {
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

func (r bookRepo) Search(ctx context.Context, searchQuery string) ([]model.Book, error) {
	termQuery := elastic.NewFuzzyQuery("title", searchQuery)
	searchResult, err := r.es.Client.Search().Index("go8-books").Query(termQuery).Pretty(true).
		Do(ctx)
	if err != nil {
		return nil, err
	}

	var book model.Book
	var books []model.Book

	for _, item := range searchResult.Each(reflect.TypeOf(book)) {
		t := item.(model.Book)
		books = append(books, t)
	}

	return books, nil
}