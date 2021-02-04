package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"testing"
	"time"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/joho/godotenv/autoload"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"

	"github.com/gmhafiz/go8/cmd/extmigrate/migrate"
	"github.com/gmhafiz/go8/configs"
	"github.com/gmhafiz/go8/internal/domain/book"
	"github.com/gmhafiz/go8/internal/models"
	"github.com/gmhafiz/go8/internal/middleware"
)

//go:generate mockgen -package mock -source ../../repository.go -destination=../../mock/mock_repository.go

var (
	repo book.Repository
)

const uniqueDBName = "postgres_test"

func TestMain(m *testing.M) {
	// must go back to project's root path to get to the .env and ./database/migrations/ folder
	to := "../../../../../"
	err := os.Chdir(to)
	if err != nil {
		log.Fatalln(err)
	}
	err = godotenv.Load(".env")
	if err != nil {
		log.Println(err)
	}
	cfg := configs.DockerTestCfg()
	cfg.Name = uniqueDBName

	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}

	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13",
		Env: []string{
			"POSTGRES_USER=" + cfg.User,
			"POSTGRES_PASSWORD=" + cfg.Pass,
			"POSTGRES_DB=" + uniqueDBName,
			"TZ=UTC",
			"PG_TZ=UTC",
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: cfg.Port},
			},
		},
	}

	resource, err := pool.RunWithOptions(&opts, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalln("error running docker container")
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Pass, uniqueDBName)
	// Docker layer network is different on Mac
	if runtime.GOOS == "darwin" {
		cfg.Host = net.JoinHostPort(resource.GetBoundIP("5432/tcp"), resource.GetPort("5432/tcp"))
	}

	if err = pool.Retry(func() error {
		db, err := sqlx.Open(cfg.Dialect, dsn)
		if err != nil {
			return err
		}
		repo = NewBookRepository(db)
		return db.Ping()
	}); err != nil {
		log.Fatalf("could not connect to docker: %s", err.Error())
	}

	migrate.Start()

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Printf("could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestBookRepository_Create(t *testing.T) {
	dt := "2020-01-01T15:04:05Z"
	timeWant, err := time.Parse(time.RFC3339, dt)
	if err != nil {
		t.Fatal(err)
	}
	bookTest := &models.Book{
		Title:         "test1",
		PublishedDate: timeWant,
		Description:   "test1",
		ImageURL: null.String{
			String: "http://example.com/image.png",
			Valid:  true,
		},
	}

	bookID, err := repo.Create(context.Background(), bookTest)

	assert.NoError(t, err)
	assert.NotEqual(t, 0, bookID)
}

func TestRepository_Find(t *testing.T) {
	dt := "2020-01-01T15:04:05Z"
	timeWant, err := time.Parse(time.RFC3339, dt)
	if err != nil {
		t.Fatal(err)
	}
	bookWant := &models.Book{
		Title:         "test2",
		PublishedDate: timeWant,
		Description:   "test2",
	}
	bookID, err := repo.Create(context.Background(), bookWant)
	if err != nil {
		t.Fatal(err)
	}

	bookGot, err := repo.Find(context.Background(), bookID)
	if err != nil {
		t.Fatal(err)
	}

	assert.Equal(t, bookGot.Title, bookWant.Title)
	assert.Equal(t, bookGot.Description, bookWant.Description)
	assert.Equal(t, bookGot.PublishedDate.String(), bookWant.PublishedDate.String())
}

func TestRepository_All(t *testing.T) {
	pagination := middleware.Pagination{
		Page:      1,
		Size:      10,
		Direction: "asc",
	}
	ctx := context.Background()
	ctx = context.WithValue(ctx, middleware.PaginationKey, pagination)

	books, err := repo.All(ctx)

	assert.NoError(t, err)
	assert.Len(t, books, 2)
}

func TestRepository_Update(t *testing.T) {
	ctx := context.Background()
	dt := "2020-01-01T15:04:05Z"
	timeWant, err := time.Parse(time.RFC3339, dt)
	if err != nil {
		assert.NoError(t, err)
	}

	want := &models.Book{
		BookID:        1,
		Title:         "updated title 1",
		PublishedDate: timeWant,
		ImageURL: null.String{
			String: "http://example.com/image2.png",
			Valid:  true,
		},
		Description: "updated description",
	}

	got, err := repo.Update(ctx, want)

	assert.NoError(t, err)
	assert.Equal(t, want.BookID, got.BookID)
	assert.Equal(t, want.Title, got.Title)
	assert.Equal(t, want.Description, got.Description)
	assert.Equal(t, want.PublishedDate.String(), got.PublishedDate.String())
	assert.Equal(t, want.ImageURL.String, got.ImageURL.String)
}

func TestRepository_Delete(t *testing.T) {
	ctx := context.Background()

	err := repo.Delete(ctx, 1)

	assert.NoError(t, err)

	got, err := repo.Find(ctx, 1)

	assert.Nil(t, got)
	assert.Error(t, err, sql.ErrNoRows)
}

func TestRepository_Search(t *testing.T) {
	ctx := context.Background()

	req := &book.Request{
		Title: "test2",
	}

	got, err := repo.Search(ctx, req)

	assert.NoError(t, err)
	assert.Len(t, got, 1)
}
