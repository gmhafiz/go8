package postgres

import (
	"context"
	"log"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"

	"github.com/gmhafiz/go8/internal/domain/book"
	"github.com/gmhafiz/go8/internal/models"
	"github.com/gmhafiz/go8/internal/utility/filter"
)

//go:generate  mockgen -package=mock -source=internal/domain/book/repository.go -destination internal/domain/book/mock/mock_repository.go

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
		Title:         "test1",
		PublishedDate: timeWant(t),
		Description:   "test1",
		ImageURL: null.String{
			String: "https://example.com/image.png",
			Valid:  true,
		},
	}
	mock.ExpectPrepare("^INSERT INTO books").
		ExpectQuery().
		WithArgs(bookTest.Title, bookTest.PublishedDate, bookTest.ImageURL, bookTest.Description).
		WillReturnRows(sqlmock.NewRows([]string{"book_id"}).AddRow(1))

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
		BookID:        1,
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
			Page:          1,
			Size:          10,
			DisablePaging: true,
			Search:        false,
		},
	}
	mock.ExpectQuery("SELECT (.+) FROM books ORDER BY").
		WillReturnRows(sqlmock.NewRows([]string{"book_id", "title", "published_date", "image_url", "description"}).
			AddRow(mockBook.BookID, mockBook.Title, mockBook.PublishedDate, mockBook.ImageURL.String, mockBook.Description),
		)

	gotBooks, err := repo.List(context.Background(), f)

	assert.NoError(t, err)
	t.Log(gotBooks)
}

func testPaginate(t *testing.T) {
	db, mock := NewMock()
	repo := New(db)

	mockBook := &models.Book{
		BookID:        1,
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
			Page:          1,
			Size:          10,
			DisablePaging: false,
			Search:        false,
		},
	}
	mock.ExpectQuery("SELECT (.+) FROM books ORDER BY").
		WithArgs(f.Base.Size, f.Base.Page).
		WillReturnRows(sqlmock.NewRows([]string{"book_id", "title", "published_date", "image_url", "description"}).
			AddRow(mockBook.BookID, mockBook.Title, mockBook.PublishedDate, mockBook.ImageURL.String, mockBook.Description),
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

	bookCols := []string{"book_id", "title", "published_date", "image_url", "description"}

	mockBook := &models.Book{
		BookID:        1,
		Title:         "test1",
		PublishedDate: timeWant(t),
		ImageURL: null.String{
			String: "https://example.com/image.png",
			Valid:  true,
		},
		Description: "test1",
	}

	mock.ExpectQuery("^SELECT (.+) FROM books where book_id").
		WithArgs(bookID).
		WillReturnRows(sqlmock.NewRows(bookCols).AddRow(mockBook.BookID, mockBook.Title, mockBook.PublishedDate, mockBook.ImageURL, mockBook.Description))

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
		BookID:        1,
		Title:         "test1",
		PublishedDate: timeWant(t),
		ImageURL: null.String{
			String: "https://example.com/image.png",
			Valid:  true,
		},
		Description: "test1",
	}

	mock.ExpectExec("UPDATE books set title").
		WithArgs(mockBook.Title, mockBook.Description, mockBook.PublishedDate, mockBook.ImageURL.String, mockBook.BookID).
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

//func TestRepository_Create(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	m := mock.NewMockRepository(ctrl)
//
//	var expectID int64
//	expectID = 1
//	var err error
//	ctx := context.Background()
//	dt := "2020-01-01T15:04:05Z"
//	timeWant, err := time.Parse(time.RFC3339, dt)
//		if err != nil {
//			t.Fatal(err)
//		}
//	bookTest := &models.Book{
//		Title:         "test1",
//		PublishedDate: timeWant,
//		Description:   "test1",
//		ImageURL: null.String{
//			String: "https://example.com/image.png",
//			Valid:  true,
//		},
//	}
//
//	m.EXPECT().Create(ctx, bookTest).Times(1).Return(expectID, err)
//
//	repo := New(nil)
//	g, err :=repo.Create(ctx, bookTest)
//	t.Log(g)
//
//	gotID, err := m.Create(ctx, bookTest)
//	if err != nil {
//		return
//	}
//
//	assert.Equal(t, expectID, gotID)
//}
//
//func TestRepository_Read(t *testing.T) {
//	ctrl := gomock.NewController(t)
//	defer ctrl.Finish()
//
//	_ = mock.NewMockRepository(ctrl)
//}

//
////go:generate mockgen -package mock -source ../../repository.go -destination=../../mock/mock_repository.go
//
//var (
//	repo book.Repository
//)
//
//const uniqueDBName = "postgres_test"
//
//func TestMain(m *testing.M) {
//	// must go back to project's root path to get to the .env and ./database/migrations/ folder
//	to := "../../../../../"
//	err := os.Chdir(to)
//	if err != nil {
//		log.Fatalln(err)
//	}
//	err = godotenv.Load(".env")
//	if err != nil {
//		log.Println(err)
//	}
//	cfg := configs.DockerTestCfg()
//	cfg.Name = uniqueDBName
//
//	pool, err := dockertest.NewPool("")
//	if err != nil {
//		log.Fatalf("could not connect to docker: %s", err)
//	}
//
//	opts := dockertest.RunOptions{
//		Repository: "postgres",
//		Tag:        "13",
//		Env: []string{
//			"POSTGRES_USER=" + cfg.User,
//			"POSTGRES_PASSWORD=" + cfg.Pass,
//			"POSTGRES_DB=" + uniqueDBName,
//			"TZ=UTC",
//			"PG_TZ=UTC",
//		},
//		ExposedPorts: []string{"5432"},
//		PortBindings: map[docker.Port][]docker.PortBinding{
//			"5432": {
//				{HostIP: "0.0.0.0", HostPort: cfg.Port},
//			},
//		},
//	}
//
//	resource, err := pool.RunWithOptions(&opts, func(config *docker.HostConfig) {
//		// set AutoRemove to true so that stopped container goes away by itself
//		config.AutoRemove = true
//		config.RestartPolicy = docker.RestartPolicy{
//			Name: "no",
//		}
//	})
//	if err != nil {
//		log.Fatalln("error running docker container")
//	}
//
//	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
//		cfg.Host, cfg.Port, cfg.User, cfg.Pass, uniqueDBName)
//	// Docker layer network is different on Mac
//	if runtime.GOOS == "darwin" {
//		cfg.Host = net.JoinHostPort(resource.GetBoundIP("5432/tcp"), resource.GetPort("5432/tcp"))
//	}
//
//	if err = pool.Retry(func() error {
//		db, err := sqlx.Open(cfg.Dialect, dsn)
//		if err != nil {
//			return err
//		}
//		repo = New(db)
//		return db.Ping()
//	}); err != nil {
//		log.Fatalf("could not connect to docker: %s", err.Error())
//	}
//
//	code := m.Run()
//
//	if err := pool.Purge(resource); err != nil {
//		log.Printf("could not purge resource: %s", err)
//	}
//
//	os.Exit(code)
//}
//
//func TestBookRepository_Create(t *testing.T) {
//	migrate.Start()
//
//	dt := "2020-01-01T15:04:05Z"
//	timeWant, err := time.Parse(time.RFC3339, dt)
//	if err != nil {
//		t.Fatal(err)
//	}
//	bookTest := &models.Book{
//		Title:         "test1",
//		PublishedDate: timeWant,
//		Description:   "test1",
//		ImageURL: null.String{
//			String: "http://example.com/image.png",
//			Valid:  true,
//		},
//	}
//
//	bookID, err := repo.Create(context.Background(), bookTest)
//
//	assert.NoError(t, err)
//	assert.NotEqual(t, 0, bookID)
//
//	migrate.Down()
//}
//
//func TestRepository_Find(t *testing.T) {
//	migrate.Start()
//
//	dt := "2020-01-01T15:04:05Z"
//	timeWant, err := time.Parse(time.RFC3339, dt)
//	if err != nil {
//		t.Fatal(err)
//	}
//	bookWant := &models.Book{
//		Title:         "test2",
//		PublishedDate: timeWant,
//		Description:   "test2",
//	}
//	bookID, err := repo.Create(context.Background(), bookWant)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	bookGot, err := repo.Read(context.Background(), bookID)
//	if err != nil {
//		t.Fatal(err)
//	}
//
//	assert.Equal(t, bookGot.Title, bookWant.Title)
//	assert.Equal(t, bookGot.Description, bookWant.Description)
//	assert.Equal(t, bookGot.PublishedDate.String(), bookWant.PublishedDate.String())
//
//	migrate.Down()
//}
//
//func TestRepository_List(t *testing.T) {
//	migrate.Start()
//
//	ctx := context.Background()
//
//	f := &book.Filter{
//		Base: filter.Filter{
//			Page:   1,
//			Size:   10,
//			Search: false,
//		},
//		Title:         "",
//		Description:   "",
//		PublishedDate: "",
//	}
//	books, err := repo.List(ctx, f)
//
//	assert.NoError(t, err)
//	assert.Len(t, books, 0)
//
//	migrate.Down()
//}
//
//func TestRepository_Update(t *testing.T) {
//	migrate.Start()
//
//	ctx := context.Background()
//	dt := "2020-01-01T15:04:05Z"
//	timeWant, err := time.Parse(time.RFC3339, dt)
//	if err != nil {
//		assert.NoError(t, err)
//	}
//
//	want := &models.Book{
//		BookID:        1,
//		Title:         "updated title 1",
//		PublishedDate: timeWant,
//		ImageURL: null.String{
//			String: "http://example.com/image2.png",
//			Valid:  true,
//		},
//		Description: "updated description",
//	}
//
//	err = repo.Update(ctx, want)
//
//	assert.NoError(t, err)
//
//	migrate.Down()
//}
//
//func TestRepository_Delete(t *testing.T) {
//	migrate.Start()
//
//	ctx := context.Background()
//
//	err := repo.Delete(ctx, 1)
//
//	assert.NoError(t, err)
//
//	got, err := repo.Read(ctx, 1)
//
//	assert.Nil(t, got)
//	assert.Error(t, err, sql.ErrNoRows)
//
//	migrate.Down()
//}
//
//func TestRepository_Search(t *testing.T) {
//	migrate.Start()
//
//	ctx := context.Background()
//
//	re := book.Filter{
//		Base:          filter.Filter{},
//		Title:         "test2",
//		Description:   "",
//		PublishedDate: "",
//	}
//
//	got, err := repo.Search(ctx, &re)
//
//	assert.NoError(t, err)
//	assert.Len(t, got, 0)
//
//	migrate.Down()
//}
