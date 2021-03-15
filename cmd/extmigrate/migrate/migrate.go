package migrate

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"github.com/gmhafiz/go8/configs"
)

const (
	databaseMigrationPath = "file://database/migrations/"
)

var (
	db *sql.DB
	m  *migrate.Migrate
)

func Start() {
	cfg := configs.New()
	Up(cfg, ".")
}

func Up(cfg *configs.Configs, changeDirTo string) {
	err := os.Chdir(changeDirTo)
	if err != nil {
		log.Fatal(err)
	}
	_, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DockerTest.Host, cfg.DockerTest.Port, cfg.DockerTest.User, cfg.DockerTest.Pass, cfg.DockerTest.Name)
	db, err = sql.Open(cfg.DockerTest.Driver, dsn)
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}

	if cfg.DockerTest.Driver == "postgres" {
		driver, err := postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			log.Fatalf("error instantiating database: %v", err)
		}
		m, err = migrate.NewWithDatabaseInstance(
			databaseMigrationPath, cfg.DockerTest.Driver, driver,
		)
		if err != nil {
			log.Fatalf("error connecting to database: %v", err)
		}

		if len(os.Args) < 2 {
			log.Fatal("usage:")
		}

		_ = m.Up()
	}
}

func Down() {
	_ = m.Down()
}
