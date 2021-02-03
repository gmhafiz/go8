package postgres

import (
	"context"
	"log"
	"time"

	"github.com/friendsofgo/errors"
	"github.com/jmoiron/sqlx"

	"github.com/gmhafiz/go8/internal/domain/book"
	"github.com/gmhafiz/go8/internal/middleware"
	"github.com/gmhafiz/go8/internal/models"
)

type repository struct {
	db *sqlx.DB
}

const (
	InsertIntoBooks         = "INSERT INTO books (title, published_date, image_url, description) VALUES ($1, $2, $3, $4) RETURNING book_id"
	SelectFromBooks         = "SELECT * FROM books ORDER BY created_at DESC"
	SelectFromBooksPaginate = "SELECT * FROM books ORDER BY created_at DESC LIMIT $1 OFFSET $2"
	SelectBookByID          = "SELECT * FROM books where book_id = $1"
	UpdateBook              = "UPDATE books set title = $1, description = $2, published_date = $3, image_url = $4, updated_at = $5 where book_id = $6"
	DeleteByID              = "DELETE FROM books where book_id = ($1)"
	SearchBooks             = "SELECT * FROM books where title like '%' || $1 || '%' and description like '%'|| $2 || '%'"
)

func NewBookRepository(db *sqlx.DB) *repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, book *models.Book) (int64, error) {
	stmt, err := r.db.PrepareContext(ctx, InsertIntoBooks)

	if err != nil {
		return 0, err
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	var bookID int64
	err = stmt.QueryRowContext(ctx, book.Title, book.PublishedDate, book.ImageURL, book.Description).Scan(&bookID)
	if err != nil {
		return 0, err
	}

	return bookID, nil
}

func (r *repository) All(ctx context.Context) ([]*models.Book, error) {
	page := ctx.Value(middleware.PaginationKey).(middleware.Pagination).Page
	size := ctx.Value(middleware.PaginationKey).(middleware.Pagination).Size

	if page == 0 && size == 0 {
		var books []*models.Book
		err := r.db.SelectContext(ctx, &books, SelectFromBooks)
		if err != nil {
			return nil, errors.Wrap(err, "error fetching books")
		}

		return books, nil

	} else {
		var books []*models.Book
		err := r.db.SelectContext(ctx, &books, SelectFromBooksPaginate, size, page*(page-1))
		if err != nil {
			return nil, errors.Wrap(err, "error fetching books")
		}

		return books, nil
	}
}

func (r *repository) Find(ctx context.Context, bookID int64) (*models.Book, error) {
	stmt, err := r.db.Prepare(SelectBookByID)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	var b models.Book
	err = r.db.GetContext(ctx, &b, SelectBookByID, bookID)
	if err != nil {
		return nil, err
	}

	return &b, err
}

func (r *repository) Update(ctx context.Context, book *models.Book) (*models.Book, error) {
	now := time.Now()

	_, err := r.db.ExecContext(ctx, UpdateBook, book.Title, book.Description,
		book.PublishedDate, book.ImageURL, now, book.BookID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	bookDB, err := r.Find(ctx, book.BookID)
	if err != nil {
		log.Println(err)
		return nil, err
	}

	return bookDB, nil
}

func (r *repository) Delete(ctx context.Context, bookID int64) error {
	_, err := r.db.ExecContext(ctx, DeleteByID, bookID)

	return err
}

func (r *repository) Search(ctx context.Context, req *book.Request) ([]*models.Book, error) {
	var books []*models.Book
	err := r.db.SelectContext(ctx, &books, SearchBooks, req.Title, req.Description)
	if err != nil {
		return nil, err
	}

	return books, nil
}

// Close attaches the provider and close the connection
func (r *repository) Close() {
	defer func() {
		err := r.db.Close()
		if err != nil {
			log.Println(err)
		}
	}()
}

// Up attaches the provider and create the table
func (r *repository) Up() error {
	ctx := context.Background()

	query := "CREATE TABLE IF NOT EXISTS books(book_id bigserial, title varchar(255) not null, published_date timestamp with time zone not null, image_url varchar(255), description text not null, created_at timestamp with time zone default current_timestamp, updated_at timestamp with time zone default current_timestamp, deleted_at timestamp with time zone, primary key (book_id))"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	_, err = stmt.ExecContext(ctx)
	return err
}

// Drop attaches the provider and drop the table
func (r *repository) Drop() error {
	ctx := context.Background()

	query := "DROP TABLE IF EXISTS books cascade"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer func() {
		err = stmt.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	_, err = stmt.ExecContext(ctx)
	return err
}