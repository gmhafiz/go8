package database

import (
	"fmt"
	"log"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/gmhafiz/go8/configs"
	"github.com/gmhafiz/go8/internal/utility/database"
)

func NewSqlx(cfg *configs.Configs) *sqlx.DB {
	dsn := fmt.Sprintf("%s://%s:%s/%s?user=%s&password=%s&sslmode=disable",
		cfg.Database.Driver,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
		cfg.Database.User,
		cfg.Database.Pass,
	)
	db, err := sqlx.Open(cfg.Database.Driver, dsn)
	if err != nil {
		log.Fatal(err)
	}

	database.Alive(db.DB)

	db.DB.SetMaxOpenConns(cfg.Database.MaxConnectionPool)

	return db
}
