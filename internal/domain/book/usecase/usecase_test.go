package usecase

import (
	"context"
	"fmt"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jinzhu/now"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/stretchr/testify/assert"
	"github.com/volatiletech/null/v8"

	"github.com/gmhafiz/go8/internal/domain/book"
	"github.com/gmhafiz/go8/internal/domain/book/repository/postgres"
	"github.com/gmhafiz/go8/internal/model"
)

var (
	repo book.Repository
)

var (
	user     = "postgres"
	password = "secret"
	dbName   = "postgres"
	port     = "5433"
	dialect  = "postgres"
	dsn      = "postgres://%s:%s@localhost:%s/%s?sslmode=disable"
)

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("could not connect to docker: %s", err)
	}

	opts := dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "13",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + dbName,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port},
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
		log.Println("error running docker container")
	}

	dsn = fmt.Sprintf(dsn, user, password, port, dbName)

	if err = pool.Retry(func() error {
		db, err := sqlx.Open(dialect, dsn)
		if err != nil {
			return err
		}
		repo = postgres.NewBookRepository(db)
		return db.Ping()
	}); err != nil {
		log.Fatalf("could not connect to docker: %s", err.Error())
	}

	defer func() {
		repo.Close()
	}()

	err = repo.Drop()
	if err != nil {
		panic(err)
	}

	err = repo.Up()
	if err != nil {
		panic(err)
	}

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Fatalf("could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestBookUseCase_Create(t *testing.T) {
	uc := NewBookUseCase(repo)

	dt := "2006-01-02 15:04:05"
	//timeWant, err := time.Parse(time.RFC3339, dt)
	//if err != nil {
	//	t.Fatal("error parsing time")
	//}
	location, err := time.LoadLocation("Australia/Sydney")
	if err != nil {
		t.Fatalf("error loading location")
	}

	myConfig := &now.Config{
		WeekStartDay: time.Sunday,
		TimeLocation: location,
		TimeFormats:  []string{"2006-01-02 15:04:05"},
	}
	timeWant, err := myConfig.Parse(dt)
	if err != nil {
		t.Fatal("error parsing time")
	}
	bookWant := &model.Book{
		Title:         "title",
		PublishedDate: timeWant,
		ImageURL: null.String{
			String: "https://example.com/image.png",
			Valid:  true,
		},
		Description: null.String{
			String: "description",
			Valid:  true,
		},
	}

	bookGot, err := uc.Create(context.Background(), bookWant.Title, bookWant.Description.String,
		bookWant.ImageURL.String, bookWant.PublishedDate.String())
	if err != nil {
		t.Fatal(err)
	}
	assert.NotEqual(t, bookGot.BookID, 0)
	assert.Equal(t, bookGot.Title, bookWant.Title)
	assert.Equal(t, bookGot.PublishedDate.String(), bookWant.PublishedDate.UTC().String())
	assert.Equal(t, bookGot.Description, bookWant.Description)
	assert.Equal(t, bookGot.ImageURL, bookWant.ImageURL)
}
