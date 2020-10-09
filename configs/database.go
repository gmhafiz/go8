package configs

import (
	"fmt"
	"os"

	"database/sql"
	_ "github.com/lib/pq"
	"github.com/rs/zerolog"
)

type Database struct {
	Driver  string
	Host    string
	Port    string
	Name    string
	User    string
	Pass    string
	SslMode string
}

func DataStore() *Database {
	return &Database{
		Driver:  os.Getenv("DB_DRIVER"),
		Host:    os.Getenv("DB_HOST"),
		Port:    os.Getenv("DB_PORT"),
		Name:    os.Getenv("DB_NAME"),
		User:    os.Getenv("DB_USER"),
		Pass:    os.Getenv("DB_PASS"),
		SslMode: os.Getenv("DB_SSL_MODE"),
	}
}

func NewDatabase(log zerolog.Logger, cfg *Configs) *sql.DB {
	dsn := fmt.Sprintf("%s://%s/%s?sslmode=%s&user=%s&password=%s",
		cfg.Database.Driver,
		cfg.Database.Host,
		cfg.Database.Name,
		cfg.Database.SslMode,
		cfg.Database.User,
		cfg.Database.Pass,
	)

	db, err := sql.Open(cfg.Database.Driver, dsn)
	if err != nil {
		log.Fatal().Msg(err.Error())
		return nil
	}

	err = db.Ping()
	if err != nil {
		log.Fatal().Msg(err.Error())
		return nil
	}

	return db
}
