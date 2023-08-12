package database

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/gmhafiz/go8/config"
	_ "github.com/gmhafiz/go8/ent/gen/runtime"
	"github.com/gmhafiz/go8/internal/utility/database"
)

func NewSqlx(cfg config.Database) *sqlx.DB {
	var dsn string
	switch cfg.Driver {
	case "postgres", "pgx":
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
			cfg.Host, cfg.Port, cfg.User, cfg.Pass, cfg.Name)
		cfg.Driver = "pgx"
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@(%s:%d)/%s?parseTime=true",
			cfg.User,
			cfg.Pass,
			cfg.Host,
			cfg.Port,
			cfg.Name,
		)
	default:
		log.Fatal("Must choose a database driver")
	}

	db, err := sqlx.Open(cfg.Driver, dsn)
	if err != nil {
		log.Fatal(err)
	}

	database.Alive(db.DB)

	return db
}
