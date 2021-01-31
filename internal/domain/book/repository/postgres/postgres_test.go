package postgres

import (
	"context"
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

	"github.com/gmhafiz/go8/cmd/extmigrate/migrate"
	"github.com/gmhafiz/go8/configs"
	"github.com/gmhafiz/go8/internal/domain/book"
	"github.com/gmhafiz/go8/internal/models"
)

var (
	repo book.Test
)

const uniqueDBName = "postgres_test"

func TestMain(m *testing.M) {
	// must go back to project's root path to get to the .env and ./database/migrations/ folder
	changeDirTo := "../../../../../"
	err := os.Chdir(changeDirTo)
	if err != nil {
		log.Fatalln(err)
	}
	err = godotenv.Load(".env")
	if err != nil {
		log.Println(err)
	}
	cfg := configs.DockerTestCfg()

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
		log.Print("error running docker container")
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

	defer func() {
		repo.Close()
	}()

	dbCfg := &configs.Configs{
		Database: &configs.Database{
			Driver:  cfg.Dialect,
			Host:    cfg.Host,
			Port:    cfg.Port,
			Name:    uniqueDBName,
			User:    cfg.User,
			Pass:    cfg.Pass,
			SslMode: cfg.SslMode,
		},
	}
	migrate.Up(dbCfg, ".")

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
		Title:         "test11",
		PublishedDate: timeWant,
		Description:   "test11",
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
		Title:         "test11",
		PublishedDate: timeWant,
		Description:   "test11",
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
