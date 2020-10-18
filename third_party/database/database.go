package database

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq"
	"github.com/rs/zerolog"

	"go8ddd/configs"
)

func New(log zerolog.Logger, cfg *configs.Configs) *sql.DB {
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

