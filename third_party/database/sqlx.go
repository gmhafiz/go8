package database

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"go.nhat.io/otelsql"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"

	"github.com/gmhafiz/go8/config"
	_ "github.com/gmhafiz/go8/ent/gen/runtime"
	"github.com/gmhafiz/go8/internal/utility/database"
)

func NewSqlx(cfg config.Database) *sqlx.DB {
	var dsn string
	switch cfg.Driver {
	case "postgres", "pgx":
		dsn = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%s",
			cfg.Host, cfg.Port, cfg.User, cfg.Pass, cfg.Name,cfg.SslMode)
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
	driverName, err := otelsql.Register(cfg.Driver,
		otelsql.AllowRoot(),
		otelsql.TraceQueryWithoutArgs(),
		otelsql.WithSystem(semconv.DBSystemPostgreSQL),
	)
	if err != nil {
		_ = fmt.Errorf("otelsql driver: %v", err)
	}
	db, err := sqlx.Open(driverName, dsn)
	if err != nil {
		log.Fatal(err)
	}

	database.Alive(db.DB)

	return db
}
