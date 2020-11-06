package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"

	"go8ddd/configs"
	"go8ddd/third_party/logger"
)

const Version = "v0.2.0"

func init() {

}

func main() {
	log := logger.New(Version)
	cfg := configs.New(log)

	source := "file://database/migrations"

	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Pass, cfg.Database.Name)
	//cmdString := fmt.Sprintf("migrate -source %s -database %s up", source, dsn)

	db, err := sql.Open(cfg.Database.Driver, dsn)
	if err != nil {
		log.Error().Msg("error opening database")
		return
	}

	if cfg.Database.Driver == "postgres" {
		driver, err := postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			log.Error().Msg("error instantiating database")
			return
		}
		m, err := migrate.NewWithDatabaseInstance(
			source,
			cfg.Database.Driver, driver,
		)

		if len(os.Args) < 2 {
			log.Error().Msg("usage:")
			return
		}

		command := os.Args[1]
		if command == "up" {
			_ = m.Up()
		}
	}

	log.Info().Msg("done.")
}
