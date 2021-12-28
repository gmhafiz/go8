package postgres

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"

	"github.com/gmhafiz/go8/internal/domain/book"
	"github.com/gmhafiz/go8/internal/models"
	"github.com/gmhafiz/go8/internal/utility/filter"
)

//go:generate  mockgen -package=mock -source=../../repository.go -destination ../../mock/mock_repository.go

func NewMock() (*sqlx.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		log.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	sqlxDB := sqlx.NewDb(db, "sqlmock")

	return sqlxDB, mock
}

func TestRepository_Create(t *testing.T) {
	db, mock := NewMock()
	repo := New(db)

	expectID := int64(1)
	ctx := context.Background()

	bookTest := &models.Book{
		ID:            1,
		Title:         "test1",
		PublishedDate: timeWant(t),
		ImageURL: null.String{
			String: "https://example.com/image.png",
			Valid:  true,
		},
		Description: "test1",
	}
	mock.ExpectExec("^INSERT INTO \"books.*").
		WithArgs(bookTest.Title, bookTest.PublishedDate, bookTest.ImageURL, bookTest.Description).
		WillReturnResult(sqlmock.NewResult(1, 1))

	gotID, err := repo.Create(ctx, bookTest)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, expectID, gotID)
}

func TestRepository_List(t *testing.T) {
	testPaginate(t)
	testDisablePaging(t)
}

func testDisablePaging(t *testing.T) {
	db, mock := NewMock()
	repo := New(db)

	mockBook := &models.Book{
		ID:            1,
		Title:         "test1",
		PublishedDate: timeWant(t),
		ImageURL: null.String{
			String: "https://example.com/image.png",
			Valid:  true,
		},
		Description: "test1",
	}
	f := &book.Filter{
		Base: filter.Filter{
			Offset:        1,
			Limit:         10,
			DisablePaging: true,
			Search:        false,
		},
	}
	mock.ExpectQuery("SELECT (.+) FROM books ORDER BY").
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "published_date", "image_url", "description"}).
			AddRow(mockBook.ID, mockBook.Title, mockBook.PublishedDate, mockBook.ImageURL.String, mockBook.Description),
		)

	gotBooks, err := repo.List(context.Background(), f)

	assert.NoError(t, err)
	t.Log(gotBooks)
}

func testPaginate(t *testing.T) {
	db, mock := NewMock()
	repo := New(db)

	mockBook := &models.Book{
		ID:            1,
		Title:         "test1",
		PublishedDate: timeWant(t),
		ImageURL: null.String{
			String: "https://example.com/image.png",
			Valid:  true,
		},
		Description: "test1",
	}
	f := &book.Filter{
		Base: filter.Filter{
			Offset:        1,
			Limit:         10,
			DisablePaging: false,
			Search:        false,
		},
	}
	mock.ExpectQuery("SELECT (.+) FROM books ORDER BY").
		WithArgs(f.Base.Limit, f.Base.Offset).
		WillReturnRows(sqlmock.NewRows([]string{"id", "title", "published_date", "image_url", "description"}).
			AddRow(mockBook.ID, mockBook.Title, mockBook.PublishedDate, mockBook.ImageURL.String, mockBook.Description),
		)

	gotBooks, err := repo.List(context.Background(), f)

	assert.NoError(t, err)
	t.Log(gotBooks)
}

func TestRepository_Read(t *testing.T) {
	db, mock := NewMock()
	repo := New(db)

	ctx := context.Background()
	bookID := int64(1)

	bookCols := []string{"id", "title", "published_date", "image_url", "description"}

	mockBook := &models.Book{
		ID:            1,
		Title:         "test1",
		PublishedDate: timeWant(t),
		ImageURL: null.String{
			String: "https://example.com/image.png",
			Valid:  true,
		},
		Description: "test1",
	}

	mock.ExpectQuery("^SELECT (.+) FROM books where id").
		WithArgs(bookID).
		WillReturnRows(sqlmock.NewRows(bookCols).AddRow(mockBook.ID, mockBook.Title, mockBook.PublishedDate, mockBook.ImageURL, mockBook.Description))

	gotBook, err := repo.Read(ctx, bookID)
	if err != nil {
		t.Fatal(err)
	}

	assert.NotNil(t, t, gotBook)
}

func TestRepository_Update(t *testing.T) {
	db, mock := NewMock()
	repo := New(db)

	mockBook := &models.Book{
		ID:            1,
		Title:         "test1",
		PublishedDate: timeWant(t),
		ImageURL: null.String{
			String: "https://example.com/image.png",
			Valid:  true,
		},
		Description: "test1",
	}

	mock.ExpectExec("UPDATE books set title").
		WithArgs(mockBook.Title, mockBook.Description, mockBook.PublishedDate, mockBook.ImageURL.String, mockBook.ID).
		WillReturnResult(sqlmock.NewErrorResult(nil))

	err := repo.Update(context.Background(), mockBook)

	assert.NoError(t, err)
}

func TestRepository_Delete(t *testing.T) {
	db, mock := NewMock()
	repo := New(db)

	expectID := int64(1)
	mock.ExpectExec("DELETE FROM books").
		WithArgs(expectID).
		WillReturnResult(sqlmock.NewResult(expectID, expectID))

	err := repo.Delete(context.Background(), expectID)

	assert.NoError(t, err)
}

func timeWant(t *testing.T) time.Time {
	t.Helper()
	dt := "2020-01-01T15:04:05Z"
	timeWant, _ := time.Parse(time.RFC3339, dt)
	return timeWant
}
