package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/lib/pq"

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

	DBAlive(db)

	return db
}

func DBAlive(db *sql.DB) {
	log.Println("connecting to database... ")
	for {
		// Ping by itself is un-reliable, the connections are cached. This
		// ensures that the database is still running by executing a harmless
		// dummy query against it.
		_, err := db.Exec("SELECT true")
		if err == nil {
			return
		}
		log.Println("retrying...")
		time.Sleep(time.Second)
	}
}
