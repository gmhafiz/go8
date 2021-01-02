package database

import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"log"

	"github.com/gmhafiz/go8/configs"
)

func New(cfg *configs.Configs) *sql.DB {
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
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Ping by itself is un-reliable, the connections are cached. This
	// ensures that the database is still running by executing a harmless
	// dummy query against it.
	if _, err = db.Exec("SELECT true"); err != nil {
		log.Fatal(err)
	}

	return db
}

func NewSqlx(cfg *configs.Configs) *sqlx.DB {
	dsn := fmt.Sprintf("%s://%s/%s?sslmode=%s&user=%s&password=%s",
		cfg.Database.Driver,
		cfg.Database.Host,
		cfg.Database.Name,
		cfg.Database.SslMode,
		cfg.Database.User,
		cfg.Database.Pass, )

	db, err := sqlx.Open(cfg.Database.Driver, dsn)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
	}

	// Ping by itself is un-reliable, the connections are cached. This
	// ensures that the database is still running by executing a harmless
	// dummy query against it.
	if _, err = db.Exec("SELECT true"); err != nil {
		log.Fatal(err)
	}

	return db
}
