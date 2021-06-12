package database

import (
	"fmt"
	"log"
	"net"
	"os"
	"runtime"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/joho/godotenv"
	_ "github.com/jackc/pgx/stdlib"
	_ "github.com/joho/godotenv/autoload"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"

	"github.com/cmd/migrate/migrate"
	"github.com/configs"
	"github.com/internal/domain/author"
)

//go:generate mockgen -package mock -source ../../repository.go -destination=../../mock/mock_repository.go

var (
	repo author.Repository
)

const uniqueDBName = "db_test"

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
		repo = New(db)
		return db.Ping()
	}); err != nil {
		log.Fatalf("could not connect to docker: %s", err.Error())
	}

	code := m.Run()

	if err := pool.Purge(resource); err != nil {
		log.Printf("could not purge resource: %s", err)
	}

	os.Exit(code)
}

func TestAuthorRepository_Create(t *testing.T) {
	migrate.Start()

	migrate.Down()
}

func TestRepository_Find(t *testing.T) {
	migrate.Start()

	migrate.Down()
}

func TestRepository_List(t *testing.T) {
	migrate.Start()

	migrate.Down()
}

func TestRepository_Update(t *testing.T) {
	migrate.Start()

	migrate.Down()
}

func TestRepository_Delete(t *testing.T) {
	migrate.Start()

	migrate.Down()
}

func TestRepository_Search(t *testing.T) {
	migrate.Start()

	migrate.Down()
}
