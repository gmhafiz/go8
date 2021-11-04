package migrate

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/stdlib"

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
	up(cfg, ".")
}

func up(cfg *configs.Configs, changeDirTo string) {
	err := os.Chdir(changeDirTo)
	if err != nil {
		log.Fatal(err)
	}
	_, err = os.Getwd()
	if err != nil {
		log.Fatal(err)
	}

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Pass, cfg.Database.Name)
	db, err = sql.Open(cfg.Database.Driver, dsn)
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}

	var driver database.Driver
	switch cfg.Database.Driver {
	case "postgres":
		d, err := postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			log.Fatalf("error instantiating database: %v", err)
		}
		driver = d
	case "mysql":
		d, err := mysql.WithInstance(db, &mysql.Config{})
		if err != nil {
			log.Fatalf("error instantiating database: %v", err)
		}
		driver = d
	}

	m, err = migrate.NewWithDatabaseInstance(
		databaseMigrationPath, cfg.Database.Driver, driver,
	)
	if err != nil {
		log.Fatalf("error connecting to database: %v", err)
	}

	if len(os.Args) < 2 {
		log.Fatal("usage:")
	}

	_ = m.Up()
}
