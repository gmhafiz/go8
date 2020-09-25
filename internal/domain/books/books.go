package books

import (
	"context"
	"database/sql"

	v "github.com/go-redis/redis/v8"
	"github.com/gocraft/work"
	p "github.com/gomodule/redigo/redis"
	"github.com/rs/zerolog"

	"eight/internal/models"
	"eight/pkg/jobs"
)

type HandlerBooks struct {
	store    store
	cache    bookCacheStore
	client   p.Conn
	enqueuer *work.Enqueuer
	jobs     *jobs.Jobs
	conn     p.Conn
}

type Context struct {
	customerID int64
}

func NewService(db *sql.DB, logger zerolog.Logger, rdb *v.Client, conn p.Conn, jobs *jobs.Jobs, enqueuer *work.Enqueuer) (*HandlerBooks, error) {
	bookStore, err := newStore(db, logger)
	if err != nil {
		return nil, err
	}

	cacheStore, err := newCacheStore(rdb, conn, logger)

	return &HandlerBooks{
		store:    bookStore,
		cache:    cacheStore,
		jobs:     jobs,
		conn:     conn,
		enqueuer: enqueuer,
	}, nil
}

func (h *HandlerBooks) AllBooks(ctx context.Context) (models.BookSlice, error) {
	booksRedis, err := h.cache.GetBooks(ctx)
	if err == nil && booksRedis != nil {
		return booksRedis, nil
	}
	books, err := h.store.All(ctx)
	if err != nil {
		return nil, err
	}

	_ = h.cache.SetBooks(ctx, &books)

	_, err = h.enqueuer.Enqueue("send_email", work.Q{"address": "test@example.com",
		"subject": "hello world", "customer_id": 4})
	if err != nil {
		return books, err
	}
	h.jobs.WorkerPool.Start()

	return books, nil
}

func (h *HandlerBooks) CreateBook(ctx context.Context, book *models.Book) (*models.Book, error) {
	return h.store.CreateBook(ctx, book)
}

func (h *HandlerBooks) GetBook(ctx context.Context, bookID int64) (*models.Book, error) {
	return h.store.GetBook(ctx, bookID)
}

func (h *HandlerBooks) Delete(ctx context.Context, bookID int64) error {
	return h.store.Delete(ctx, bookID)
}

func (h *HandlerBooks) Ping() error {
	return h.store.Ping()
}
