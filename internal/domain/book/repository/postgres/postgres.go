package postgres

import (
	"context"

	"github.com/friendsofgo/errors"
	"github.com/jmoiron/sqlx"

	"github.com/gmhafiz/go8/internal/middleware"
	"github.com/gmhafiz/go8/internal/model"
	"github.com/gmhafiz/go8/internal/resource"
)

type repository struct {
	db *sqlx.DB
}

func NewBookRepository(db *sqlx.DB) *repository {
	return &repository{db: db}
}

func (r *repository) Create(ctx context.Context, book *model.Book) (int64, error) {
	query := "INSERT INTO books (title, published_date, image_url, description) VALUES ($1, $2, $3, $4) RETURNING book_id"
	stmt, err := r.db.PrepareContext(ctx, query)

	if err != nil {
		return 0, err
	}
	defer stmt.Close()

	var bookID int64
	err = stmt.QueryRowContext(ctx, book.Title, book.PublishedDate, book.ImageURL, book.Description).Scan(&bookID)
	if err != nil {
		return 0, err
	}

	return bookID, nil
}

func (r *repository) All(ctx context.Context) ([]resource.BookDB, error) {
	page := ctx.Value(middleware.PaginationKey).(middleware.Pagination).Page
	size := ctx.Value(middleware.PaginationKey).(middleware.Pagination).Size

	if page == 0 && size == 0 {
		query := "SELECT * FROM books ORDER BY created_at DESC"
		rows, err := r.db.QueryContext(ctx, query)
		if err != nil {
			return nil, errors.Wrap(err, "error fetching books")
		}
		var books []resource.BookDB
		for rows.Next() {
			var book resource.BookDB
			err := rows.Scan(&book.BookID, &book.Title, &book.PublishedDate,
				&book.ImageURL, &book.Description, &book.CreatedAt, &book.UpdatedAt, &book.DeletedAt)
			if err != nil {
				return nil, errors.Wrap(err, "error scanning book")
			}
			books = append(books, book)
		}
		return books, nil

	} else {
		query := "SELECT * FROM books ORDER BY created_at DESC LIMIT $1 OFFSET $2 "
		rows, err := r.db.QueryContext(ctx, query, size, page*(page-1))
		if err != nil {
			return nil, errors.Wrap(err, "error fetching books")
		}
		var books []resource.BookDB
		for rows.Next() {
			var book resource.BookDB
			err = rows.Scan(&book.BookID, &book.Title, &book.PublishedDate,
				&book.ImageURL, &book.Description, &book.CreatedAt, &book.UpdatedAt, &book.DeletedAt)
			if err != nil {
				return nil, errors.Wrap(err, "error scanning book")
			}
			books = append(books, book)
		}
		return books, nil
	}
}

func (r *repository) Find(ctx context.Context, bookID int64) (*model.Book, error) {
	query := "SELECT * FROM books where book_id = $1"
	stmt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	bookDB := resource.BookDB{}
	err = stmt.QueryRow(bookID).Scan(&bookDB.BookID, &bookDB.Title, &bookDB.PublishedDate,
		&bookDB.ImageURL, &bookDB.Description, &bookDB.CreatedAt, &bookDB.UpdatedAt,
		&bookDB.DeletedAt)
	if err != nil {
		return nil, err
	}

	b := &model.Book{
		BookID:        bookDB.BookID,
		Title:         bookDB.Title,
		PublishedDate: bookDB.PublishedDate,
		ImageURL:      bookDB.ImageURL,
		Description:   bookDB.Description,
		CreatedAt:     bookDB.CreatedAt,
		UpdatedAt:     bookDB.UpdatedAt,
		DeletedAt:     bookDB.DeletedAt,
	}

	return b, err
}

// Close attaches the provider and close the connection
func (r *repository) Close() {
	r.db.Close()
}

// Up attaches the provider and create the table
func (r *repository) Up() error {
	ctx := context.Background()

	query := "CREATE TABLE IF NOT EXISTS books(book_id bigserial, title varchar(255) not null, published_date timestamp with time zone not null, image_url varchar(255), description text not null, created_at timestamp with time zone default current_timestamp, updated_at timestamp with time zone default current_timestamp, deleted_at timestamp with time zone, primary key (book_id))"
	stmt, err := r.db.PrepareContext(ctx, query)
	if err != nil {
		return err
	}
	defer stmt.Close()

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
	defer stmt.Close()

	_, err = stmt.ExecContext(ctx)
	return err
}
