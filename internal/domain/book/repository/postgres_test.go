package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/gmhafiz/go8/internal/utility/message"
	"log"
	"math"
	"os"
	"testing"
	"time"

	_ "github.com/gmhafiz/go8/ent/gen/runtime"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
	"github.com/ory/dockertest/v3"
	"github.com/ory/dockertest/v3/docker"
	"github.com/stretchr/testify/assert"

	"github.com/gmhafiz/go8/internal/domain/book"
	"github.com/gmhafiz/go8/internal/utility/filter"
)

var (
	dockerDB *db
)

var (
	startTime = time.Now()
)

type db struct {
	Conn *sql.DB
	Sqlx *sqlx.DB
}

func (d db) migrate(upOrDown string) {
	driver, err := postgres.WithInstance(d.Conn, &postgres.Config{})
	if err != nil {
		log.Fatal(err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://../../../../database/migrations/",
		"postgres",
		driver,
	)
	if err != nil {
		log.Fatal(err)
	}

	switch upOrDown {
	case "down":
		if err = m.Down(); err != nil {
			log.Fatalln(err)
		}
	default:
		if err = m.Up(); err != nil {
			log.Fatalln(err)
		}
	}
}

func TestMain(m *testing.M) {
	// uses a sensible default on windows (tcp/http) and linux/osx (socket)
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// pulls an image, creates a container based on it and runs it
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "15",
		Env: []string{
			"POSTGRES_PASSWORD=secret",
			"POSTGRES_USER=user_name",
			"POSTGRES_DB=dbname",
			"listen_addresses = '*'",
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{Name: "no"}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	hostAndPort := resource.GetHostPort("5432/tcp")
	databaseUrl := fmt.Sprintf("postgres://user_name:secret@%s/dbname?sslmode=disable", hostAndPort)
	log.Println(databaseUrl)

	_ = resource.Expire(120) // Tell docker to hard kill the container in 120 seconds

	dockerDB = &db{}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	pool.MaxWait = 120 * time.Second
	if err = pool.Retry(func() error {
		dockerDB.Conn, err = sql.Open("postgres", databaseUrl)
		if err != nil {
			return err
		}
		return dockerDB.Conn.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	// Performing a migration this way means all tests in this package shares
	// the same db schema across all unit test.
	// If isolation is needed, then do away with using `testing.M`. Do a
	// migration for each test handler instead.
	dockerDB.migrate("up")

	// We can access database with m.hostAndPort or m.databaseUrl
	// port changes everytime a new docker instance is run
	code := m.Run()

	// You can't defer this because os.Exit doesn't care for defer
	if err := pool.Purge(resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestRepository_Create(t *testing.T) {
	type args struct {
		ctx context.Context
		req *book.CreateRequest
	}
	type want struct {
		lastInsertID int
		err          error
	}

	type test struct {
		name string
		args
		want
	}

	tests := []test{
		{
			name: "simple",
			args: args{
				ctx: context.Background(),
				req: &book.CreateRequest{
					Title:         "title",
					PublishedDate: "2020-01-01T15:04:05Z",
					ImageURL:      "https://example.com/image.png",
					Description:   "description",
				},
			},
			want: want{
				lastInsertID: 1,
				err:          nil,
			},
		},
		{
			name: "adding a second book should return ID=2",
			args: args{
				ctx: context.Background(),
				req: &book.CreateRequest{
					Title:         "2",
					PublishedDate: "2020-01-01T15:04:05Z",
					ImageURL:      "https://example.com/image.png",
					Description:   "description",
				},
			},
			want: want{
				lastInsertID: 2,
				err:          nil,
			},
		},
		{
			name: "empty strings",
			args: args{
				ctx: context.Background(),
				req: &book.CreateRequest{
					Title:         "",
					PublishedDate: "",
					ImageURL:      "",
					Description:   "",
				},
			},
			want: want{
				lastInsertID: 0,
				err:          errors.New("repository.Book.Create"),
			},
		},
	}

	client := dbClient()
	repo := New(client)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := repo.Create(test.args.ctx, test.args.req)
			assert.Equal(t, test.want.err, err)
			assert.Equal(t, test.want.lastInsertID, got)
		})
	}
}

func TestRepository_List(t *testing.T) {
	type args struct {
		ctx context.Context
		f   *book.Filter
	}
	type want struct {
		books []*book.Schema
		err   error
	}

	type test struct {
		name string
		args
		want
	}

	//startTime := time.Now()

	timeParsed, err := time.Parse(time.RFC3339, "2020-01-01T15:04:05Z")
	assert.Nil(t, err)

	tests := []test{
		{
			name: "Should return one",
			args: args{
				ctx: context.Background(),
				f: &book.Filter{
					Base: filter.Filter{
						Page:          0,
						Offset:        0,
						Limit:         1,
						DisablePaging: false,
						Sort:          nil,
						Search:        false,
					},
					Title:         "",
					Description:   "",
					PublishedDate: "",
				},
			},
			want: want{
				books: []*book.Schema{
					{
						ID:            2,
						Title:         "2",
						PublishedDate: timeParsed,
						ImageURL:      "https://example.com/image.png",
						Description:   "description",
						CreatedAt:     time.Now(),
						UpdatedAt:     time.Now(),
						DeletedAt:     sql.NullTime{},
					},
				},
				err: nil,
			},
		},
		{
			name: "should return all",
			args: args{
				ctx: context.Background(),
				f: &book.Filter{
					Base: filter.Filter{
						Page:          0,
						Offset:        0,
						Limit:         10,
						DisablePaging: false,
						Sort:          nil,
						Search:        false,
					},
					Title:         "",
					Description:   "",
					PublishedDate: "",
				},
			},
			want: want{
				books: []*book.Schema{
					{
						ID:            2,
						Title:         "2",
						PublishedDate: timeParsed,
						ImageURL:      "https://example.com/image.png",
						Description:   "description",
					},
					{
						ID:            1,
						Title:         "title",
						PublishedDate: timeParsed,
						ImageURL:      "https://example.com/image.png",
						Description:   "description",
					},
				},
				err: nil,
			},
		},
		{
			name: "disable paging",
			args: args{
				ctx: context.Background(),
				f: &book.Filter{
					Base: filter.Filter{
						Page:          0,
						Offset:        0,
						Limit:         10,
						DisablePaging: true,
						Sort:          nil,
						Search:        false,
					},
					Title:         "",
					Description:   "",
					PublishedDate: "",
				},
			},
			want: want{
				books: []*book.Schema{
					{
						ID:            2,
						Title:         "2",
						PublishedDate: timeParsed,
						ImageURL:      "https://example.com/image.png",
						Description:   "description",
					},
					{
						ID:            1,
						Title:         "title",
						PublishedDate: timeParsed,
						ImageURL:      "https://example.com/image.png",
						Description:   "description",
					},
				},
				err: nil,
			},
		},
		{
			name: "filter cannot be nil",
			args: args{
				ctx: context.Background(),
				f:   nil,
			},
			want: want{
				books: nil,
				err:   errors.New("filter cannot be nil"),
			},
		},
	}

	client := dbClient()
	repo := New(client)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := repo.List(test.args.ctx, test.args.f)
			assert.Equal(t, test.want.err, err)

			if err != nil {
				assert.Nil(t, test.want.books)
				return
			}

			for i, val := range got {
				assert.Equal(t, test.want.books[i].ID, val.ID)
				assert.Equal(t, test.want.books[i].Title, val.Title)
				assert.Equal(t, test.want.books[i].PublishedDate.UTC(), val.PublishedDate.UTC())
				assert.Equal(t, test.want.books[i].ImageURL, val.ImageURL)
				assert.Equal(t, test.want.books[i].Description, val.Description)
				assert.True(t, startTime.Before(got[i].CreatedAt) || startTime.Equal(got[i].CreatedAt))
				assert.True(t, startTime.Before(got[i].UpdatedAt) || startTime.Equal(got[i].UpdatedAt))
				assert.Equal(t, test.want.books[i].DeletedAt, got[i].DeletedAt)
			}
		})
	}
}

func TestRepository_Read(t *testing.T) {
	type args struct {
		context.Context
		int
	}
	type want struct {
		book *book.Schema
		err  error
	}
	type test struct {
		name string
		args
		want
	}

	//startTime := time.Now()

	timeParsed, err := time.Parse(time.RFC3339, "2020-02-17T00:00:00Z")
	assert.Nil(t, err)

	createOneBook := &book.CreateRequest{
		Title:         "title",
		PublishedDate: "2020-02-17T00:00:00Z",
		ImageURL:      "https://example.com/image.png",
		Description:   "description",
	}

	tests := []test{
		{
			name: "simple",
			args: args{
				Context: context.Background(),
				int:     3,
			},
			want: want{
				book: &book.Schema{
					ID:            3,
					Title:         "title",
					PublishedDate: timeParsed,
					ImageURL:      "https://example.com/image.png",
					Description:   "description",
					CreatedAt:     time.Now(),
					UpdatedAt:     time.Now(),
					DeletedAt:     sql.NullTime{},
				},
				err: nil,
			},
		},
		{
			name: "simulate error",
			args: args{
				Context: context.Background(),
				int:     -0,
			},
			want: want{
				book: nil,
				err:  message.ErrBadRequest,
			},
		},
	}

	client := dbClient()
	repo := New(client)

	_, err = repo.Create(context.Background(), createOneBook)
	assert.Nil(t, err)

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := repo.Read(test.args.Context, test.args.int)
			assert.Equal(t, test.want.err, err)
			if err != nil {
				assert.Nil(t, test.want.book)
				return
			}

			assert.Equal(t, test.want.book.ID, got.ID)
			assert.Equal(t, test.want.book.Title, got.Title)
			assert.Equal(t, test.want.book.PublishedDate.UTC(), got.PublishedDate.UTC())
			assert.Equal(t, test.want.book.ImageURL, got.ImageURL)
			assert.Equal(t, test.want.book.Description, got.Description)
			assert.True(t, startTime.Before(got.CreatedAt) || startTime.Equal(got.CreatedAt))
			assert.True(t, startTime.Before(got.UpdatedAt) || startTime.Equal(got.UpdatedAt))
			assert.Equal(t, test.want.book.DeletedAt, got.DeletedAt)

		})
	}
}

func TestRepository_Update(t *testing.T) {
	type args struct {
		context.Context
		*book.UpdateRequest
	}
	type want struct {
		err  error
		book *book.Schema
	}
	type test struct {
		name string
		args
		want
	}

	timeParsed, err := time.Parse(time.RFC3339, "2020-02-17T00:00:00Z")
	assert.Nil(t, err)

	tests := []test{
		{
			name: "update title",
			args: args{
				Context: context.Background(),
				UpdateRequest: &book.UpdateRequest{
					ID:            3,
					Title:         "updated title",
					PublishedDate: "2020-02-17T00:00:00Z",
					ImageURL:      "https://example.com/image.png",
					Description:   "description",
				},
			},
			want: want{
				err: nil,
				book: &book.Schema{
					ID:            3,
					Title:         "updated title",
					PublishedDate: timeParsed,
					ImageURL:      "https://example.com/image.png",
					Description:   "description",
					CreatedAt:     time.Time{}, // will be before time.Now()
					UpdatedAt:     time.Now(),
					DeletedAt:     sql.NullTime{},
				},
			},
		},
		{
			name: "invalid ID",
			args: args{
				Context: context.Background(),
				UpdateRequest: &book.UpdateRequest{
					ID:            -1,
					Title:         "updated title",
					PublishedDate: "2020-02-17T00:00:00Z",
					ImageURL:      "https://example.com/image.png",
					Description:   "description",
				},
			},
			want: want{
				book: &book.Schema{},
				err:  sql.ErrNoRows,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := dbClient()
			repo := New(client)

			err := repo.Update(test.args.Context, test.args.UpdateRequest)
			assert.Equal(t, test.want.err, err)

			if err != nil {
				return
			}

			got, err := repo.Read(test.args.Context, test.args.UpdateRequest.ID)
			assert.Nil(t, err)

			assert.Equal(t, test.want.book.ID, got.ID)
			assert.Equal(t, test.want.book.Title, got.Title)
			assert.Equal(t, test.want.book.PublishedDate.UTC(), got.PublishedDate.UTC())
			assert.Equal(t, test.want.book.Description, got.Description)
			assert.Equal(t, test.want.book.ImageURL, got.ImageURL)
			assert.True(t, startTime.Before(got.CreatedAt) || startTime.Equal(got.CreatedAt))
			assert.True(t, startTime.Before(got.UpdatedAt) || startTime.Equal(got.UpdatedAt))
			assert.Equal(t, test.want.book.DeletedAt, got.DeletedAt)

			//assert.True(t, test.want.book.DeletedAt.Before(*got.DeletedAt))
		})
	}
}

func TestRepository_Delete(t *testing.T) {
	type args struct {
		context.Context
		bookID int
	}
	type want struct {
		err error
	}
	type test struct {
		name string
		args
		want
	}

	tests := []test{
		{
			name: "simple",
			args: args{
				Context: context.Background(),
				bookID:  1,
			},
			want: want{
				err: nil,
			},
		},
		{
			name: "delete non-existent ID",
			args: args{
				Context: context.Background(),
				bookID:  math.MaxInt - 1,
			},
			want: want{
				err: errors.New("ID not found: sql: no rows in result set"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := dbClient()
			repo := New(client)

			err := repo.Delete(test.args.Context, test.args.bookID)

			if err != nil {
				assert.Equal(t, test.want.err.Error(), err.Error())
				return
			}
			assert.Equal(t, test.want.err, err)

		})
	}
}

func TestRepository_Search(t *testing.T) {
	type args struct {
		context.Context
		f *book.Filter
	}
	type want struct {
		books []*book.Schema
		err   error
	}
	type test struct {
		name string
		args
		want
	}

	timeParsed, err := time.Parse(time.RFC3339, "2020-01-01T15:04:05Z")
	assert.Nil(t, err)

	tests := []test{
		{
			name: "finds ID=2",
			args: args{
				Context: context.Background(),
				f: &book.Filter{
					Base: filter.Filter{
						Page:          0,
						Offset:        0,
						Limit:         10,
						DisablePaging: false,
						Sort:          nil,
						Search:        true,
					},
					Title:         "2",
					Description:   "",
					PublishedDate: "",
				},
			},
			want: want{
				books: []*book.Schema{
					{
						ID:            2,
						Title:         "2",
						PublishedDate: timeParsed,
						ImageURL:      "https://example.com/image.png",
						Description:   "description",
					},
				},
				err: nil,
			},
		},
		{
			name: "nil filter",
			args: args{
				Context: context.Background(),
				f:       nil,
			},
			want: want{
				books: nil,
				err:   errors.New("filter cannot be nil"),
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			client := dbClient()
			repo := New(client)

			got, err := repo.Search(test.args.Context, test.args.f)
			assert.Equal(t, test.want.err, err)

			if err != nil {
				return
			}

			for i, val := range got {
				assert.Equal(t, test.want.books[i].ID, val.ID)
				assert.Equal(t, test.want.books[i].Title, val.Title)
				assert.Equal(t, test.want.books[i].PublishedDate.UTC(), val.PublishedDate.UTC())
				assert.Equal(t, test.want.books[i].ImageURL, val.ImageURL)
				assert.Equal(t, test.want.books[i].Description, val.Description)
				assert.True(t, startTime.Before(got[i].CreatedAt) || startTime.Equal(got[i].CreatedAt))
				assert.True(t, startTime.Before(got[i].UpdatedAt) || startTime.Equal(got[i].UpdatedAt))
				assert.Equal(t, test.want.books[i].DeletedAt, got[i].DeletedAt)
			}

		})
	}
}

func dbClient() *sqlx.DB {
	sqlxDB := sqlx.NewDb(dockerDB.Conn, "pgx")
	dockerDB.Sqlx = sqlxDB
	return sqlxDB
}
