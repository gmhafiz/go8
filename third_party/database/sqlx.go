package database

import (
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/gmhafiz/go8/configs"
)

func NewSqlx(cfg *configs.Configs) *sqlx.DB {
	dsn := fmt.Sprintf("%s://%s/%s?sslmode=%s&user=%s&password=%s",
		cfg.Database.Driver,
		cfg.Database.Host,
		cfg.Database.Name,
		cfg.Database.SslMode,
		cfg.Database.User,
		cfg.Database.Pass)

	db, err := sqlx.Open(cfg.Database.Driver, dsn)
	if err != nil {
		log.Fatal(err)
	}

	DBAlive(db.DB)

	return db
}
