package database

import (
	"database/sql"
	"log"
	"time"
)

func Alive(db *sql.DB) {
	log.Println("Connecting to database... ")
	for {
		// Ping by itself is un-reliable, the connections are cached. This
		// ensures that the database is still running by executing a harmless
		// dummy query against it.
		_, err := db.Exec("SELECT true")
		if err == nil {
			log.Println("Database connected")
			return
		}
		log.Println(err)
		log.Println("retrying...")
		time.Sleep(1 * time.Second)
	}
}
