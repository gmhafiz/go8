package database

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"

	"github.com/gmhafiz/go8/config"
	_ "github.com/gmhafiz/go8/ent/gen/runtime"
	"github.com/gmhafiz/go8/internal/utility/database"
)

func NewSqlx(cfg *config.Config) *sqlx.DB {
	var dsn string
	switch cfg.Database.Driver {
	case "postgres", "pgx":
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			cfg.Database.Host, cfg.Database.Port, cfg.Database.User, cfg.Database.Pass, cfg.Database.Name)
	case "mysql":
		dsn = fmt.Sprintf("%s:%s@(%s:%s)/%s?parseTime=true",
			cfg.Database.User,
			cfg.Database.Pass,
			cfg.Database.Host,
			cfg.Database.Port,
			cfg.Database.Name,
		)
	default:
		log.Fatal("Must choose a database driver")
	}

	db, err := sqlx.Open(cfg.Database.Driver, dsn)
	if err != nil {
		log.Fatal(err)
	}

	database.Alive(db.DB)

	db.DB.SetMaxOpenConns(cfg.Database.MaxConnectionPool)

	return db
}
